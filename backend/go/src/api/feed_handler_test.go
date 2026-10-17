package api

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
	"github.com/gin-gonic/gin"
)

type MockLatestNewsRepository struct {
	news []models.News
	err  error
}

func (m *MockLatestNewsRepository) GetLatestNews(limit int) ([]models.News, error) {
	return m.news, m.err
}

func serveRSS(t *testing.T, repo LatestNewsRepositoryInterface) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	h := &FeedHandler{newsRepo: repo, siteURL: "https://example.com"}
	r.GET("/news/rss", h.GetRSS)
	r.ServeHTTP(w, httptest.NewRequest("GET", "/news/rss", nil))
	return w
}

func TestGetRSS_Success(t *testing.T) {
	createdAt := time.Date(2026, 10, 17, 3, 10, 0, 0, time.UTC)
	withSummary := models.News{
		ID:             42,
		Title:          "Go 1.27 & <generics>",
		Description:    "raw description",
		SummarizedText: "要約テキスト",
		Category:       "Programming",
		ArticleUrl:     "https://example.com/a?x=1&y=2",
		CreatedAt:      createdAt,
	}
	withoutSummary := models.News{
		ID:          41,
		Title:       "No summary",
		Description: "fallback description",
		ArticleUrl:  "https://example.com/b",
		CreatedAt:   createdAt.Add(-time.Hour),
	}

	w := serveRSS(t, &MockLatestNewsRepository{news: []models.News{withSummary, withoutSummary}})

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/rss+xml; charset=utf-8" {
		t.Errorf("unexpected Content-Type: %s", ct)
	}

	var feed rss
	if err := xml.Unmarshal(w.Body.Bytes(), &feed); err != nil {
		t.Fatalf("failed to parse feed: %v\n%s", err, w.Body.String())
	}
	if feed.Version != "2.0" || feed.Channel.Link != "https://example.com" {
		t.Errorf("unexpected channel: %+v", feed.Channel)
	}
	if feed.Channel.LastBuildDate != createdAt.Format(time.RFC1123Z) {
		t.Errorf("unexpected lastBuildDate: %s", feed.Channel.LastBuildDate)
	}
	if len(feed.Channel.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(feed.Channel.Items))
	}

	first := feed.Channel.Items[0]
	if first.Title != withSummary.Title || first.Link != withSummary.ArticleUrl {
		t.Errorf("unexpected item: %+v", first)
	}
	if first.Description != "要約テキスト" {
		t.Errorf("expected summarized text, got %q", first.Description)
	}
	if first.GUID.Value != "birdseye-news-42" || first.GUID.IsPermaLink {
		t.Errorf("unexpected guid: %+v", first.GUID)
	}
	if first.PubDate != createdAt.Format(time.RFC1123Z) {
		t.Errorf("unexpected pubDate: %s", first.PubDate)
	}
	if feed.Channel.Items[1].Description != "fallback description" {
		t.Errorf("expected description fallback, got %q", feed.Channel.Items[1].Description)
	}
}

func TestGetRSS_Empty(t *testing.T) {
	w := serveRSS(t, &MockLatestNewsRepository{})

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var feed rss
	if err := xml.Unmarshal(w.Body.Bytes(), &feed); err != nil {
		t.Fatalf("failed to parse feed: %v", err)
	}
	if len(feed.Channel.Items) != 0 || feed.Channel.LastBuildDate != "" {
		t.Errorf("expected empty channel, got %+v", feed.Channel)
	}
}

func TestGetRSS_RepositoryError(t *testing.T) {
	w := serveRSS(t, &MockLatestNewsRepository{err: errors.New("db down")})

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}
