package news

import (
	"encoding/json"
	"fmt"

	"github.com/birdseyeapi/birds_eye_v3/go/src/ai"
	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
	"github.com/birdseyeapi/birds_eye_v3/go/src/scraping/doc"
)

const (
	ZennSourceName      = "Zenn"
	ZennBaseURL         = "https://zenn.dev"
	ZennArticleSelector = "#tech-trend > div > div > div > article > div > a"
)

var MaxArticles = 15

type ScrapeNewsByZenn struct {
	summarizer ai.Summarizer
}

type zennNextData struct {
	Props struct {
		PageProps struct {
			DailyTechArticles []struct {
				Title *string `json:"title"`
				Path  *string `json:"path"`
			} `json:"dailyTechArticles"`
		} `json:"pageProps"`
	} `json:"props"`
}

func NewScrapeNewsByZenn(summarizer ai.Summarizer) *ScrapeNewsByZenn {
	return &ScrapeNewsByZenn{
		summarizer: summarizer,
	}
}

func (s *ScrapeNewsByZenn) GetSourceBy() string {
	return ZennSourceName
}

func (s *ScrapeNewsByZenn) ExtractNews() ([]models.News, error) {
	d, err := doc.GetWebDoc(ZennBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Zennは <script id="__NEXT_DATA__"> にJSON文字列で記事データが埋め込まれている
	script := d.Find("script#__NEXT_DATA__").First()
	var nextData zennNextData
	if script.Length() > 0 {
		if err := json.Unmarshal([]byte(script.Text()), &nextData); err != nil {
			fmt.Printf("Failed to parse JSON: %v\n", err)
		}
	}

	articles := nextData.Props.PageProps.DailyTechArticles
	if len(articles) == 0 {
		return nil, fmt.Errorf("no articles found in __NEXT_DATA__")
	}
	if len(articles) > MaxArticles {
		articles = articles[:MaxArticles]
	}

	var links []articleLink
	for _, article := range articles {
		if article.Title == nil || article.Path == nil {
			continue
		}
		links = append(links, articleLink{title: *article.Title, url: ZennBaseURL + *article.Path})
	}

	return buildNewsList(s.summarizer, ZennSourceName, ZennBaseURL, links), nil
}
