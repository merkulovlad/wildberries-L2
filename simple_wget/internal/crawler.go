package internal

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
)

type Crawler struct {
	opts       Options
	downloader *Downloader
	startHost  string

	mu      sync.Mutex
	visited map[string]string // canonical URL -> local path

	sem chan struct{}
	wg  sync.WaitGroup
}

func NewCrawler(opts Options, start *url.URL) *Crawler {
	return &Crawler{
		opts:       opts,
		downloader: NewDownloader(opts.Timeout, opts.UserAgent),
		startHost:  start.Host,
		visited:    make(map[string]string),
		sem:        make(chan struct{}, opts.Concurrency),
	}
}

// Run downloads `start` and recursively follows links up to opts.Depth.
func (c *Crawler) Run(start *url.URL) {
	c.enqueue(start, 0)
	c.wg.Wait()
}

func (c *Crawler) enqueue(u *url.URL, depth int) {
	canon := canonicalize(u)

	c.mu.Lock()
	if _, ok := c.visited[canon]; ok {
		c.mu.Unlock()
		return
	}
	localPath := LocalPath(u)
	c.visited[canon] = localPath
	c.mu.Unlock()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.sem <- struct{}{}
		defer func() { <-c.sem }()
		c.fetch(u, localPath, depth)
	}()
}

func (c *Crawler) fetch(u *url.URL, localPath string, depth int) {
	log.Printf("GET %s (depth=%d)", u.String(), depth)

	resp, err := c.downloader.Get(u.String())
	if err != nil {
		log.Printf("error fetching %s: %v", u.String(), err)
		return
	}

	fullPath := filepath.Join(c.opts.OutputDir, localPath)

	if !IsHTML(resp.ContentType) {
		if err := SaveFile(fullPath, resp.Body); err != nil {
			log.Printf("error saving %s: %v", fullPath, err)
		}
		return
	}

	doc, refs, err := ExtractRefs(bytes.NewReader(resp.Body))
	if err != nil {
		log.Printf("error parsing %s: %v", u.String(), err)
		return
	}

	for _, ref := range refs {
		c.handleRef(ref, u, localPath, depth)
	}

	var buf bytes.Buffer
	if err := Render(&buf, doc); err != nil {
		log.Printf("error rendering %s: %v", u.String(), err)
		return
	}
	if err := SaveFile(fullPath, buf.Bytes()); err != nil {
		log.Printf("error saving %s: %v", fullPath, err)
	}
}

func (c *Crawler) handleRef(ref *Reference, page *url.URL, pageLocalPath string, depth int) {
	if ref.IsSrcset {
		entries := SplitSrcset(ref.Raw)
		for i := range entries {
			if local := c.resolveAndQueue(entries[i].URL, page, pageLocalPath, depth, ref.TagName); local != "" {
				entries[i].URL = local
			}
		}
		ref.Set(JoinSrcset(entries))
		return
	}

	if local := c.resolveAndQueue(ref.Raw, page, pageLocalPath, depth, ref.TagName); local != "" {
		ref.Set(local)
	}
}

// resolveAndQueue resolves raw against page, decides whether to download it,
// and returns the rewritten local link (or "" to leave the attribute alone).
func (c *Crawler) resolveAndQueue(raw string, page *url.URL, pageLocalPath string, depth int, tag string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.HasPrefix(raw, "#") ||
		strings.HasPrefix(raw, "javascript:") ||
		strings.HasPrefix(raw, "mailto:") ||
		strings.HasPrefix(raw, "data:") ||
		strings.HasPrefix(raw, "tel:") {
		return ""
	}

	target, err := page.Parse(raw)
	if err != nil {
		return ""
	}
	if target.Scheme != "http" && target.Scheme != "https" {
		return ""
	}

	// <a> links follow depth limits; assets (img/script/link/...) are always
	// pulled when they live on the same host so the page renders offline.
	isPageLink := tag == "a" || tag == "iframe"

	if c.opts.SameHost && target.Host != c.startHost {
		return ""
	}

	if isPageLink && depth >= c.opts.Depth {
		return ""
	}

	nextDepth := depth
	if isPageLink {
		nextDepth = depth + 1
	}

	c.enqueue(stripFragment(target), nextDepth)

	resourceLocal := filepath.Join(c.opts.OutputDir, LocalPath(target))
	pageOnDisk := filepath.Join(c.opts.OutputDir, pageLocalPath)
	link := RelativeLink(pageOnDisk, resourceLocal)
	if target.Fragment != "" {
		link += "#" + target.Fragment
	}
	return link
}

func canonicalize(u *url.URL) string {
	cp := *u
	cp.Fragment = ""
	return cp.String()
}

func stripFragment(u *url.URL) *url.URL {
	cp := *u
	cp.Fragment = ""
	return &cp
}

// Summary prints a short report after Run finishes.
func (c *Crawler) Summary() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return fmt.Sprintf("downloaded %d unique URLs into %s", len(c.visited), c.opts.OutputDir)
}
