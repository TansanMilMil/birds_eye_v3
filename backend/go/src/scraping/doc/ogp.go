package doc

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExtractOGImage(d *goquery.Document, pageURL string) string {
	if d == nil {
		return ""
	}

	content := ""
	d.Find("meta[property='og:image'], meta[property='og:image:url'], meta[name='og:image']").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		content = strings.TrimSpace(s.AttrOr("content", ""))
		return content == ""
	})
	if content == "" {
		return ""
	}

	ref, err := url.Parse(content)
	if err != nil {
		return ""
	}
	if ref.IsAbs() {
		return ref.String()
	}

	base, err := url.Parse(pageURL)
	if err != nil {
		return ""
	}
	return base.ResolveReference(ref).String()
}
