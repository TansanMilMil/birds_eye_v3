package reaction

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
	"github.com/birdseyeapi/birds_eye_v3/go/src/scraping/doc"
	"github.com/tebeka/selenium"
)

const (
	TwitterSourceName    = "Twitter"
	TwitterBaseURL       = "https://search.yahoo.co.jp/realtime/search"
	TwitterTweetSelector = `p[class^="Tweet_body__"]`
)

// TwitterMaxReactions caps how many tweets are kept per article, matching the
// limit used for the Hatena bookmark reactions.
var TwitterMaxReactions = 10

type ScrapeReactionsByTwitter struct{}

func NewScrapeReactionsByTwitter() *ScrapeReactionsByTwitter {
	return &ScrapeReactionsByTwitter{}
}

func (s *ScrapeReactionsByTwitter) GetSourceBy() string {
	return TwitterSourceName
}

func (s *ScrapeReactionsByTwitter) ExtractReactions(driver selenium.WebDriver, newsId uint, articleURL string, title string) ([]models.NewsReaction, error) {
	var reactions []models.NewsReaction

	searchURL := TwitterBaseURL + "?ei=UTF-8&p=" + url.QueryEscape(articleURL)

	d, err := doc.GetWebDoc(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Yahoo realtime search: %v", err)
	}

	tweets := d.Find(TwitterTweetSelector)

	limit := TwitterMaxReactions
	if tweets.Length() < limit {
		limit = tweets.Length()
	}
	tweets = tweets.Slice(0, limit)

	tweets.Each(func(i int, tweet *goquery.Selection) {
		text := strings.TrimSpace(tweet.Text())
		if text == "" || text == title {
			return
		}

		reaction := models.NewsReaction{
			NewsID:          newsId,
			Author:          "twitter user",
			Comment:         text,
			ScrapedDateTime: time.Now().UTC(),
			CommentUrl:      searchURL,
		}

		reactions = append(reactions, reaction)
	})

	return reactions, nil
}
