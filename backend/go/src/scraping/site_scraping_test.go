package scraping

import (
	"errors"
	"testing"

	"github.com/birdseyeapi/birds_eye_v3/go/src/ai"
	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
)

type fakeCategorizer struct {
	results map[string]ai.CategoryResult
}

func (f *fakeCategorizer) Categorize(title, summary, sourceBy string) (ai.CategoryResult, error) {
	r, ok := f.results[title]
	if !ok {
		return ai.CategoryResult{}, errors.New("categorize failed")
	}
	return r, nil
}

func TestCategorizeNews(t *testing.T) {
	s := &SiteScraping{categorizer: &fakeCategorizer{results: map[string]ai.CategoryResult{
		"ok": {Category: ai.CategorySecurity, Confidence: 0.8},
	}}}
	newsList := []models.News{{Title: "ok"}, {Title: "ng"}}

	s.categorizeNews(newsList)

	if newsList[0].Category != ai.CategorySecurity || newsList[0].CategoryConfidence != 0.8 {
		t.Errorf("unexpected category for categorized news: %+v", newsList[0])
	}
	if newsList[1].Category != "" {
		t.Errorf("Expected empty category when categorization fails, got %q", newsList[1].Category)
	}
}

func TestCategorizeNews_NilCategorizer(t *testing.T) {
	s := &SiteScraping{}
	newsList := []models.News{{Title: "a"}}

	s.categorizeNews(newsList)

	if newsList[0].Category != "" {
		t.Errorf("Expected empty category, got %q", newsList[0].Category)
	}
}
