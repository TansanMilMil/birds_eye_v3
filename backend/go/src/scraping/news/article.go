package news

import (
	"fmt"
	"time"

	"github.com/birdseyeapi/birds_eye_v3/go/src/ai"
	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
	"github.com/birdseyeapi/birds_eye_v3/go/src/scraping/doc"
)

type articleLink struct {
	title    string
	url      string
	imageUrl string
}

func buildNews(summarizer ai.Summarizer, sourceBy, scrapedUrl string, link articleLink) (models.News, bool) {
	newsItem := models.News{
		Title:           link.title,
		SourceBy:        sourceBy,
		ScrapedUrl:      scrapedUrl,
		ScrapedDateTime: time.Now(),
		ArticleUrl:      link.url,
		ArticleImageUrl: link.imageUrl,
	}

	artDoc, err := doc.GetWebDoc(link.url)
	if err != nil {
		fmt.Printf("Failed to parse article HTML: %v\n", err)
		return models.News{}, false
	}

	if artDoc != nil && summarizer != nil {
		summary, err := summarizer.Summarize(artDoc.Text())
		if err != nil {
			fmt.Printf("Failed to summarize article: %v\n", err)
		} else {
			newsItem.SummarizedText = summary
		}
	}

	fmt.Print(".")
	return newsItem, true
}

func buildNewsList(summarizer ai.Summarizer, sourceBy, scrapedUrl string, links []articleLink) []models.News {
	var news []models.News
	for _, link := range links {
		if newsItem, ok := buildNews(summarizer, sourceBy, scrapedUrl, link); ok {
			news = append(news, newsItem)
		}
	}
	return news
}
