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
	CloudWatchSourceName      = "CloudWatch by Impress"
	CloudWatchBaseURL         = "https://cloud.watch.impress.co.jp"
	CloudWatchArticleSelector = "li.item.news"
)

var CloudWatchMaxArticles = 5

type ScrapeNewsByCloudWatchImpress struct {
	summarizer ai.Summarizer
}

func NewScrapeNewsByCloudWatchImpress(summarizer ai.Summarizer) *ScrapeNewsByCloudWatchImpress {
	return &ScrapeNewsByCloudWatchImpress{
		summarizer: summarizer,
	}
}

func (s *ScrapeNewsByCloudWatchImpress) GetSourceBy() string {
	return CloudWatchSourceName
}

func (s *ScrapeNewsByCloudWatchImpress) ExtractNews() ([]models.News, error) {
	d, err := doc.GetWebDoc(CloudWatchBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	articles := d.Find(CloudWatchArticleSelector)
	if articles.Length() == 0 {
		return nil, fmt.Errorf("no articles found with selector '%s'", CloudWatchArticleSelector)
	}

	var links []articleLink
	articles.Slice(0, CloudWatchMaxArticles).Each(func(i int, art *goquery.Selection) {
		titleElement := art.Find("p.title > a")
		artUrl := strings.TrimSpace(titleElement.AttrOr("href", ""))
		if artUrl != "" && !strings.HasPrefix(artUrl, "http") {
			artUrl = CloudWatchBaseURL + artUrl
		}

		links = append(links, articleLink{
			title: strings.TrimSpace(titleElement.Text()),
			url:   artUrl,
		})
	})

	return buildNewsList(s.summarizer, CloudWatchSourceName, CloudWatchBaseURL, links), nil
}
