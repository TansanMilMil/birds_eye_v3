package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
	"github.com/gin-gonic/gin"
)

type MockNewsRepository struct {
	news []models.News
}

func (m *MockNewsRepository) GetNews(t time.Time, c *gin.Context) []models.News {
	var result []models.News
	targetDate := t.Format("2006-01-02")
	for _, n := range m.news {
		newsDate := n.CreatedAt.Format("2006-01-02")
		if newsDate == targetDate {
			result = append(result, n)
		}
	}
	return result
}

func setupTestHandler(mockNews []models.News) *NewsHandler {
	handler := &NewsHandler{
		newsRepo: &MockNewsRepository{news: mockNews},
	}
	return handler
}

func createMockNews(date time.Time, count int) []models.News {
	var news []models.News
	for i := 0; i < count; i++ {
		n := models.News{
			Title:           "Test News",
			Description:     "Test Description",
			SourceBy:        "Test Source",
			ScrapedUrl:      "https://example.com",
			ScrapedDateTime: date,
			ArticleUrl:      "https://example.com/article",
		}
		n.CreatedAt = date
		news = append(news, n)
	}
	return news
}

func TestGetNewsByDate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	today := time.Now()
	mockNews := createMockNews(today, 5)
	yesterday := today.AddDate(0, 0, -1)
	mockNews = append(mockNews, createMockNews(yesterday, 3)...)

	handler := setupTestHandler(mockNews)

	r.GET("/news/:target_date", handler.GetNewsByDate)
	c.Request = httptest.NewRequest("GET", "/news/"+today.Format("2006-01-02"), nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string][]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	newsList, ok := resp["news"]
	if !ok {
		t.Fatal("Response did not contain 'news' key")
	}

	if len(newsList) != 5 {
		t.Errorf("Expected 5 news items, got %d", len(newsList))
	}
}

func TestGetNewsByDate_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	handler := setupTestHandler(nil)

	r.GET("/news/:target_date", handler.GetNewsByDate)
	c.Request = httptest.NewRequest("GET", "/news/invalid-date", nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid date format, got %d", w.Code)
	}
}

func TestGetNewsByDate_NoNewsFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	mockNews := createMockNews(yesterday, 3)

	handler := setupTestHandler(mockNews)

	r.GET("/news/:target_date", handler.GetNewsByDate)
	c.Request = httptest.NewRequest("GET", "/news/"+today.Format("2006-01-02"), nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string][]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	newsList, ok := resp["news"]
	if !ok {
		t.Fatal("Response did not contain 'news' key")
	}

	if len(newsList) != 0 {
		t.Errorf("Expected 0 news items, got %d", len(newsList))
	}
}

func TestGetNewsByDate_SortedByImportance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	today := time.Now()
	mockNews := createMockNews(today, 4)
	for i, importance := range []float64{0.2, 0.9, 0, 0.5} {
		mockNews[i].Importance = importance
	}

	handler := setupTestHandler(mockNews)

	r.GET("/news/:target_date", handler.GetNewsByDate)
	c.Request = httptest.NewRequest("GET", "/news/"+today.Format("2006-01-02"), nil)
	r.ServeHTTP(w, c.Request)

	var resp struct {
		News []models.GetAllNewsResponse `json:"news"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	var got []float64
	for _, n := range resp.News {
		got = append(got, n.Importance)
	}
	want := []float64{0.9, 0.5, 0.2, 0}
	if len(got) != len(want) {
		t.Fatalf("Expected %d news items, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Expected importance order %v, got %v", want, got)
		}
	}
}

type capturingNewsRepository struct {
	target time.Time
}

func (m *capturingNewsRepository) GetNews(t time.Time, c *gin.Context) []models.News {
	m.target = t
	return nil
}

func TestGetNewsByDate_TargetDateFormat(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "offset date-time", path: "/news/2026-10-14T00:00:00-04:00", want: "2026-10-14T00:00:00-04:00"},
		{name: "percent-encoded offset date-time", path: "/news/2026-10-14T00%3A00%3A00%2B09%3A00", want: "2026-10-14T00:00:00+09:00"},
		{name: "UTC date-time", path: "/news/2026-10-14T00:00:00Z", want: "2026-10-14T00:00:00Z"},
		{name: "legacy date-only defaults to JST", path: "/news/2026-10-14", want: "2026-10-14T00:00:00+09:00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			repo := &capturingNewsRepository{}
			handler := &NewsHandler{newsRepo: repo}
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			r.GET("/news/:target_date", handler.GetNewsByDate)
			r.ServeHTTP(w, httptest.NewRequest("GET", tt.path, nil))

			if w.Code != http.StatusOK {
				t.Fatalf("Expected status 200, got %d", w.Code)
			}
			if got := repo.target.Format(time.RFC3339); got != tt.want {
				t.Errorf("Expected target %s, got %s", tt.want, got)
			}
		})
	}
}

func TestGetNewsByDate_DateTimeWithoutOffset(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &NewsHandler{newsRepo: &capturingNewsRepository{}}
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/news/:target_date", handler.GetNewsByDate)
	r.ServeHTTP(w, httptest.NewRequest("GET", "/news/2026-10-14T00:00:00", nil))

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
