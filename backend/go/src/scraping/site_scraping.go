package scraping

import (
	"fmt"
	"net/url"

	"github.com/birdseyeapi/birds_eye_v3/go/src/ai"
	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
	"github.com/birdseyeapi/birds_eye_v3/go/src/scraping/news"
	"github.com/birdseyeapi/birds_eye_v3/go/src/scraping/reaction"
	"github.com/tebeka/selenium"
)

type SiteScraping struct {
	scrapers         []news.ScrapingNews
	reactionScrapers []reaction.ScrapingReaction
	categorizer      ai.Categorizer
}

func NewSiteScraping() *SiteScraping {
	summarizer := ai.NewOpenAISummarizer()

	return &SiteScraping{
		scrapers: []news.ScrapingNews{
			news.NewScrapeNewsByCloudWatchImpress(summarizer),
			news.NewScrapeNewsByHatena(summarizer),
			news.NewScrapeNewsByZenn(summarizer),
			news.NewScrapeNewsByZDNet(summarizer),
		},
		categorizer: ai.NewJevCategorizer(),
		reactionScrapers: []reaction.ScrapingReaction{
			reaction.NewScrapeReactionsByHatena(),
			reaction.NewScrapeReactionsByTwitter(),
		},
	}
}

func (s *SiteScraping) ScrapeNews() ([]models.News, error) {
	allNews := []models.News{}

	for _, scraper := range s.scrapers {
		fmt.Print(scraper.GetSourceBy() + ": scraping")
		news, err := safeExtractNews(scraper)
		if err != nil {
			fmt.Printf(" -> Error scraping from %s: %v\n", scraper.GetSourceBy(), err)
			continue
		}

		fmt.Printf(" -> scraped article: %s: %d\n", scraper.GetSourceBy(), len(news))
		allNews = append(allNews, news...)
	}

	s.categorizeNews(allNews)

	return allNews, nil
}

func (s *SiteScraping) categorizeNews(newsList []models.News) {
	if s.categorizer == nil {
		return
	}

	for i := range newsList {
		n := &newsList[i]
		result, err := s.categorizer.Categorize(n.Title, n.SummarizedText, n.SourceBy)
		if err != nil {
			fmt.Printf("Failed to categorize article %q: %v\n", n.Title, err)
			continue
		}
		n.Category = result.Category
		n.CategoryConfidence = result.Confidence
	}
}

// NewReactionDriver creates a single Selenium session to be shared across all
// reaction scrapes in one run. The caller owns its lifecycle and must Quit it.
func (s *SiteScraping) NewReactionDriver() (selenium.WebDriver, error) {
	return reaction.NewFirefoxDriver()
}

func (s *SiteScraping) ScrapeReactions(driver selenium.WebDriver, news models.News) ([]models.NewsReaction, error) {
	allReactions := []models.NewsReaction{}

	_, err := url.Parse(news.ArticleUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid article URL: %v", err)
	}

	for _, scraper := range s.reactionScrapers {
		reactions, err := safeExtractReactions(scraper, driver, news.ID, news.ArticleUrl, news.Title)
		if err != nil {
			fmt.Printf("Error scraping reactions from %s: %v\n", scraper.GetSourceBy(), err)
			continue
		}

		fmt.Printf(" -> scraped reactions: %s: %d\n", scraper.GetSourceBy(), len(reactions))
		allReactions = append(allReactions, reactions...)
	}

	return allReactions, nil
}

// safeExtractNews runs a single news scraper, converting any panic into an
// error so that a runtime failure in one site does not abort the scraping of
// the remaining sites.
func safeExtractNews(scraper news.ScrapingNews) (result []models.News, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = fmt.Errorf("panic while scraping %s: %v", scraper.GetSourceBy(), r)
		}
	}()
	return scraper.ExtractNews()
}

// safeExtractReactions runs a single reaction scraper, converting any panic
// into an error so one site's runtime failure does not abort the others.
func safeExtractReactions(scraper reaction.ScrapingReaction, driver selenium.WebDriver, newsID uint, articleURL, title string) (result []models.NewsReaction, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = fmt.Errorf("panic while scraping reactions from %s: %v", scraper.GetSourceBy(), r)
		}
	}()
	return scraper.ExtractReactions(driver, newsID, articleURL, title)
}
