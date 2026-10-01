package news

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birds_eye_v3/go/src/ai"
	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
	"github.com/birdseyeapi/birds_eye_v3/go/src/scraping/doc"
)

const (
	HatenaSourceName      = "Hatena"
	HatenaBaseURL         = "https://b.hatena.ne.jp/hotentry/it"
	HatenaArticleSelector = "#container .entrylist-contents-main"
)

var HatenaMaxArticles = 15

type ScrapeNewsByHatena struct {
	summarizer ai.Summarizer
}

func NewScrapeNewsByHatena(summarizer ai.Summarizer) *ScrapeNewsByHatena {
	return &ScrapeNewsByHatena{
		summarizer: summarizer,
	}
}

func (s *ScrapeNewsByHatena) GetSourceBy() string {
	return HatenaSourceName
}

func (s *ScrapeNewsByHatena) ExtractNews() ([]models.News, error) {
	d, err := doc.GetWebDoc(HatenaBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	articles := d.Find(HatenaArticleSelector)
	if articles.Length() == 0 {
		return nil, fmt.Errorf("no articles found with selector '%s'", HatenaArticleSelector)
	}

	var links []articleLink
	articles.Slice(0, HatenaMaxArticles).Each(func(i int, art *goquery.Selection) {
		titleElement := art.Find(".entrylist-contents-title > a")
		links = append(links, articleLink{
			title: titleElement.Text(),
			url:   titleElement.AttrOr("href", ""),
		})
	})

	return buildNewsList(s.summarizer, HatenaSourceName, HatenaBaseURL, links), nil
}
