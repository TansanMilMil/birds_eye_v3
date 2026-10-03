package doc

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestExtractOGImage(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{"absolute", `<head><meta property="og:image" content="https://example.com/a.png"></head>`, "https://example.com/a.png"},
		{"relative path", `<head><meta property="og:image" content="/img/a.png"></head>`, "https://example.com/img/a.png"},
		{"protocol relative", `<head><meta property="og:image" content="//cdn.example.com/a.png"></head>`, "https://cdn.example.com/a.png"},
		{"og:image:url fallback", `<head><meta property="og:image:url" content="https://example.com/b.png"></head>`, "https://example.com/b.png"},
		{"empty content skipped", `<head><meta property="og:image" content=""><meta property="og:image" content="https://example.com/c.png"></head>`, "https://example.com/c.png"},
		{"none", `<head><meta property="og:title" content="t"></head>`, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			if err != nil {
				t.Fatal(err)
			}
			if got := ExtractOGImage(d, "https://example.com/post/1"); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractOGImageNilDoc(t *testing.T) {
	if got := ExtractOGImage(nil, "https://example.com"); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}
