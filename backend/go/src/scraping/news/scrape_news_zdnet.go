package news

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birds_eye_v3/go/src/ai"
	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
	"github.com/birdseyeapi/birds_eye_v3/go/src/scraping/doc"
)

const (
	ZDNetSourceName      = "ZDNet Japan"
	ZDNetBaseURL         = "https://japan.zdnet.com"
	ZDNetArticleSelector = "#page-wrap > div.pg-container-main > main > section:nth-child(1) > div > ul > li"
)

var ZDNetMaxArticles = 15

type ScrapeNewsByZDNet struct {
	summarizer ai.Summarizer
}

func NewScrapeNewsByZDNet(summarizer ai.Summarizer) *ScrapeNewsByZDNet {
	return &ScrapeNewsByZDNet{
		summarizer: summarizer,
	}
}

func (s *ScrapeNewsByZDNet) GetSourceBy() string {
	return ZDNetSourceName
}

func (s *ScrapeNewsByZDNet) ExtractNews() ([]models.News, error) {
	d, err := doc.GetWebDoc(ZDNetBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}
	if d == nil {
		return nil, fmt.Errorf("failed to get document: document is nil")
	}

	articles := d.Find(ZDNetArticleSelector)
	if articles.Length() == 0 {
		return nil, fmt.Errorf("no articles found with selector '%s'", ZDNetArticleSelector)
	}
	if articles.Length() > ZDNetMaxArticles {
		articles = articles.Slice(0, ZDNetMaxArticles)
	}

	var links []articleLink
	articles.Each(func(i int, art *goquery.Selection) {
		if link, ok := parseZDNetArticle(art); ok {
			links = append(links, link)
		}
	})

	return buildNewsList(s.summarizer, ZDNetSourceName, ZDNetBaseURL, links), nil
}

func parseZDNetArticle(art *goquery.Selection) (articleLink, bool) {
	titleElement := art.Find("a > div.txt > p.txt-ttl")
	if titleElement.Length() == 0 {
		fmt.Println("Warning: Title element not found")
		return articleLink{}, false
	}

	title := strings.TrimSpace(titleElement.Text())
	if title == "" {
		fmt.Println("Warning: Empty title found")
		return articleLink{}, false
	}

	artUrlElem := art.Find("a")
	if artUrlElem.Length() == 0 {
		fmt.Println("Warning: Link element not found")
		return articleLink{}, false
	}

	artUrl := ZDNetBaseURL + strings.TrimSpace(artUrlElem.AttrOr("href", ""))
	if artUrl == ZDNetBaseURL {
		fmt.Println("Warning: Invalid article URL")
		return articleLink{}, false
	}

	imageUrl := ""
	if src, exists := art.Find("a > div.thumb > img").Attr("src"); exists {
		imageUrl = ZDNetBaseURL + src
	}

	return articleLink{title: title, url: artUrl, imageUrl: imageUrl}, true
}
