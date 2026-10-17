package api

import (
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
	repo "github.com/birdseyeapi/birds_eye_v3/go/src/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const feedItemLimit = 50

var siteURL = os.Getenv("BIRDSEYE_SITE_URL")

type LatestNewsRepositoryInterface interface {
	GetLatestNews(limit int) ([]models.News, error)
}

type FeedHandler struct {
	newsRepo LatestNewsRepositoryInterface
	siteURL  string
}

func NewFeedHandler(db *gorm.DB) *FeedHandler {
	return &FeedHandler{
		newsRepo: repo.NewNewsRepository(db),
		siteURL:  siteURL,
	}
}

type rss struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate,omitempty"`
	Items         []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	Description string  `xml:"description"`
	Category    string  `xml:"category,omitempty"`
	GUID        rssGUID `xml:"guid"`
	PubDate     string  `xml:"pubDate"`
}

type rssGUID struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

func (h *FeedHandler) GetRSS(c *gin.Context) {
	news, err := h.newsRepo.GetLatestNews(feedItemLimit)
	if err != nil {
		log.Printf("Error getting latest news for feed: %v", err)
		c.String(http.StatusInternalServerError, "failed to build feed")
		return
	}

	body, err := xml.Marshal(h.buildRSS(news))
	if err != nil {
		log.Printf("Error marshaling feed: %v", err)
		c.String(http.StatusInternalServerError, "failed to build feed")
		return
	}

	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", append([]byte(xml.Header), body...))
}

func (h *FeedHandler) buildRSS(news []models.News) rss {
	channel := rssChannel{
		Title:       "Birds Eye",
		Link:        h.siteURL,
		Description: "IT ニュースの AI 要約",
		Language:    "ja",
		Items:       make([]rssItem, len(news)),
	}
	if len(news) > 0 {
		channel.LastBuildDate = news[0].CreatedAt.Format(time.RFC1123Z)
	}

	for i, n := range news {
		description := n.SummarizedText
		if description == "" {
			description = n.Description
		}
		channel.Items[i] = rssItem{
			Title:       n.Title,
			Link:        n.ArticleUrl,
			Description: description,
			Category:    n.Category,
			// Slack's RSS app dedupes by guid, so it must stay stable per article.
			GUID:    rssGUID{Value: "birdseye-news-" + strconv.FormatUint(uint64(n.ID), 10)},
			PubDate: n.CreatedAt.Format(time.RFC1123Z),
		}
	}

	return rss{Version: "2.0", Channel: channel}
}
