package scraping

import (
	"errors"
	"reflect"
	"testing"

	"github.com/birdseyeapi/birds_eye_v3/go/src/ai"
	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
)

type fakeArticleAnalyzer struct {
	results map[string]ai.ArticleAnalysis
}

func (f *fakeArticleAnalyzer) Analyze(title, summary, sourceBy string) (ai.ArticleAnalysis, error) {
	r, ok := f.results[title]
	if !ok {
		return ai.ArticleAnalysis{}, errors.New("analyze failed")
	}
	return r, nil
}

type fakeSentimentAnalyzer struct {
	counts   ai.SentimentCounts
	err      error
	received []string
}

func (f *fakeSentimentAnalyzer) Analyze(articleTitle string, comments []string) (ai.SentimentCounts, error) {
	f.received = comments
	return f.counts, f.err
}

func TestAnalyzeNews(t *testing.T) {
	s := &SiteScraping{articleAnalyzer: &fakeArticleAnalyzer{results: map[string]ai.ArticleAnalysis{
		"ok": {Category: ai.CategorySecurity, CategoryConfidence: 0.8, Importance: 0.7},
	}}}
	newsList := []models.News{{Title: "ok"}, {Title: "ng"}}

	got := s.analyzeNews(newsList)

	if got[0].Category != ai.CategorySecurity || got[0].CategoryConfidence != 0.8 || got[0].Importance != 0.7 {
		t.Errorf("unexpected analysis for analyzed news: %+v", got[0])
	}
	if got[1].Category != "" || got[1].Importance != 0 {
		t.Errorf("Expected empty analysis when analysis fails, got %+v", got[1])
	}
	if newsList[0].Category != "" || newsList[0].Importance != 0 {
		t.Errorf("Expected input to remain unmodified, got %+v", newsList[0])
	}
}

func TestAnalyzeNews_NilAnalyzer(t *testing.T) {
	s := &SiteScraping{}
	newsList := []models.News{{Title: "a"}}

	got := s.analyzeNews(newsList)

	if got[0].Category != "" {
		t.Errorf("Expected empty category, got %q", got[0].Category)
	}
}

func TestAnalyzeReactionSentiment(t *testing.T) {
	fake := &fakeSentimentAnalyzer{counts: ai.SentimentCounts{Positive: 2, Neutral: 1, Negative: 3}}
	s := &SiteScraping{sentimentAnalyzer: fake}
	n := models.News{Title: "t", Reactions: []models.NewsReaction{{Comment: "a"}, {Comment: "b"}}}

	got := s.AnalyzeReactionSentiment(n)

	want := models.ReactionSentiment{Positive: 2, Neutral: 1, Negative: 3}
	if got != want {
		t.Errorf("Expected %+v, got %+v", want, got)
	}
	if n.ReactionSentiment != (models.ReactionSentiment{}) {
		t.Errorf("Expected input to remain unmodified, got %+v", n.ReactionSentiment)
	}
	if !reflect.DeepEqual(fake.received, []string{"a", "b"}) {
		t.Errorf("unexpected comments passed to analyzer: %v", fake.received)
	}
}

func TestAnalyzeReactionSentiment_SkipsWhenNotAnalyzable(t *testing.T) {
	cases := map[string]struct {
		analyzer  ai.ReactionSentimentAnalyzer
		reactions []models.NewsReaction
	}{
		"nil analyzer":   {nil, []models.NewsReaction{{Comment: "a"}}},
		"no reactions":   {&fakeSentimentAnalyzer{counts: ai.SentimentCounts{Positive: 1}}, nil},
		"analyzer error": {&fakeSentimentAnalyzer{counts: ai.SentimentCounts{Positive: 1}, err: errors.New("x")}, []models.NewsReaction{{Comment: "a"}}},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			s := &SiteScraping{sentimentAnalyzer: tc.analyzer}
			n := models.News{Title: "t", Reactions: tc.reactions}

			got := s.AnalyzeReactionSentiment(n)

			if got != (models.ReactionSentiment{}) {
				t.Errorf("Expected empty sentiment, got %+v", got)
			}
		})
	}
}
