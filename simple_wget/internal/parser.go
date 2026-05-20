package internal

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// linkAttrs lists the (tag, attribute) pairs we treat as resource references.
var linkAttrs = map[string][]string{
	"a":      {"href"},
	"link":   {"href"},
	"script": {"src"},
	"img":    {"src", "srcset"},
	"source": {"src", "srcset"},
	"iframe": {"src"},
	"video":  {"src", "poster"},
	"audio":  {"src"},
}

// Reference is a single link found in an HTML document.
// Rewriter callbacks mutate the underlying attribute in place.
type Reference struct {
	Raw      string
	Attr     string
	TagName  string
	IsSrcset bool
	set      func(string)
}

func (r *Reference) Set(v string) { r.set(v) }

// ExtractRefs parses HTML and returns the parsed root plus all link references.
// The original document is preserved so it can be re-serialized after rewriting.
func ExtractRefs(r io.Reader) (*html.Node, []*Reference, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, nil, err
	}

	var refs []*Reference
	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if attrs, ok := linkAttrs[n.Data]; ok {
				for i := range n.Attr {
					attr := &n.Attr[i]
					if !contains(attrs, attr.Key) {
						continue
					}
					refs = append(refs, &Reference{
						Raw:      attr.Val,
						Attr:     attr.Key,
						TagName:  n.Data,
						IsSrcset: attr.Key == "srcset",
						set:      makeSetter(attr),
					})
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	return doc, refs, nil
}

func makeSetter(attr *html.Attribute) func(string) {
	return func(v string) { attr.Val = v }
}

// Render serializes a parsed document back to HTML.
func Render(w io.Writer, doc *html.Node) error {
	return html.Render(w, doc)
}

// SplitSrcset parses a srcset attribute into URL/descriptor pairs.
func SplitSrcset(s string) []SrcsetEntry {
	var out []SrcsetEntry
	for _, part := range strings.Split(s, ",") {
		fields := strings.Fields(strings.TrimSpace(part))
		if len(fields) == 0 {
			continue
		}
		e := SrcsetEntry{URL: fields[0]}
		if len(fields) > 1 {
			e.Descriptor = strings.Join(fields[1:], " ")
		}
		out = append(out, e)
	}
	return out
}

type SrcsetEntry struct {
	URL        string
	Descriptor string
}

func JoinSrcset(entries []SrcsetEntry) string {
	parts := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.Descriptor != "" {
			parts = append(parts, e.URL+" "+e.Descriptor)
		} else {
			parts = append(parts, e.URL)
		}
	}
	return strings.Join(parts, ", ")
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
