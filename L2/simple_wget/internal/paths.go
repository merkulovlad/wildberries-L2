package internal

import (
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

// LocalPath converts a URL to a relative on-disk path under the output directory.
// Directory-style URLs (ending with "/" or with no extension) become index.html
// so the resulting tree mirrors what a browser would request.
func LocalPath(u *url.URL) string {
	host := u.Host
	p := u.Path
	if p == "" || strings.HasSuffix(p, "/") {
		p += "index.html"
	} else if ext := path.Ext(p); ext == "" {
		p += "/index.html"
	}

	if u.RawQuery != "" {
		p += "_" + sanitizeQuery(u.RawQuery)
		if path.Ext(p) == "" {
			p += ".html"
		}
	}

	parts := strings.Split(p, "/")
	for i, part := range parts {
		parts[i] = sanitizeSegment(part)
	}

	return filepath.Join(host, filepath.Join(parts...))
}

func sanitizeSegment(s string) string {
	replacer := strings.NewReplacer(
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(s)
}

func sanitizeQuery(q string) string {
	r := strings.NewReplacer("&", "_", "=", "-", "/", "_")
	return r.Replace(q)
}

// RelativeLink returns a path suitable for rewriting an href/src inside `fromPage`
// so that it points to `toResource` on disk.
func RelativeLink(fromPage, toResource string) string {
	rel, err := filepath.Rel(filepath.Dir(fromPage), toResource)
	if err != nil {
		return toResource
	}
	return filepath.ToSlash(rel)
}
