package internal

import (
	"net/url"
	"path/filepath"
	"testing"
)

func TestLocalPath(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"https://example.com/", filepath.Join("example.com", "index.html")},
		{"https://example.com/about", filepath.Join("example.com", "about", "index.html")},
		{"https://example.com/style.css", filepath.Join("example.com", "style.css")},
		{"https://example.com/img/logo.png", filepath.Join("example.com", "img", "logo.png")},
		{"https://example.com/page?id=1", filepath.Join("example.com", "page", "index.html_id-1")},
	}
	for _, c := range cases {
		u, err := url.Parse(c.in)
		if err != nil {
			t.Fatal(err)
		}
		got := LocalPath(u)
		if got != c.want {
			t.Errorf("LocalPath(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestRelativeLink(t *testing.T) {
	got := RelativeLink(
		filepath.Join("site", "example.com", "about", "index.html"),
		filepath.Join("site", "example.com", "style.css"),
	)
	want := "../style.css"
	if got != want {
		t.Errorf("RelativeLink = %q, want %q", got, want)
	}
}
