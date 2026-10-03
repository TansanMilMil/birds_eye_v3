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

type fakeOpinionClassifier struct {
	opinions map[string]bool
	err      error
	received []string
}

func (f *fakeOpinionClassifier) HasOpinion(articleTitle string, comments []string) ([]bool, error) {
	f.received = comments
	if f.err != nil {
		return nil, f.err
	}
	result := make([]bool, len(comments))
	for i, c := range comments {
		result[i] = f.opinions[c]
	}
	return result, nil
}

func reactionComments(reactions []models.NewsReaction) []string {
	comments := []string{}
	for _, r := range reactions {
		comments = append(comments, r.Comment)
	}
	return comments
}

func TestFilterOpinionReactions(t *testing.T) {
	fake := &fakeOpinionClassifier{opinions: map[string]bool{
		"これは便利そう https://example.com": true,
		"記事を読んだ":                      false,
	}}
	s := &SiteScraping{opinionClassifier: fake}
	reactions := []models.NewsReaction{
		{Comment: "https://example.com/a"},
		{Comment: "タイトル https://example.com/a"},
		{Comment: "これは便利そう https://example.com"},
		{Comment: "記事を読んだ"},
	}

	got := s.filterOpinionReactions("タイトル", reactions)

	if want := []string{"これは便利そう https://example.com"}; !reflect.DeepEqual(reactionComments(got), want) {
		t.Errorf("Expected %v, got %v", want, reactionComments(got))
	}
	if want := []string{"これは便利そう https://example.com", "記事を読んだ"}; !reflect.DeepEqual(fake.received, want) {
		t.Errorf("Expected only URL/title-free comments sent to classifier, got %v", fake.received)
	}
}

func TestFilterOpinionReactions_KeepsCandidatesWhenNotClassifiable(t *testing.T) {
	cases := map[string]ai.ReactionOpinionClassifier{
		"nil classifier":   nil,
		"classifier error": &fakeOpinionClassifier{err: errors.New("x")},
	}

	for name, classifier := range cases {
		t.Run(name, func(t *testing.T) {
			s := &SiteScraping{opinionClassifier: classifier}
			reactions := []models.NewsReaction{{Comment: "https://example.com"}, {Comment: "良い"}}

			got := s.filterOpinionReactions("t", reactions)

			if want := []string{"良い"}; !reflect.DeepEqual(reactionComments(got), want) {
				t.Errorf("Expected %v, got %v", want, reactionComments(got))
			}
		})
	}
}

func TestIsShareOnly(t *testing.T) {
	cases := map[string]struct {
		comment string
		want    bool
	}{
		"url only":           {"https://example.com/a", true},
		"title and url":      {"タイトル https://example.com/a", true},
		"title with spaces":  {"  タイトル  ", true},
		"opinion with url":   {"便利そう https://example.com/a", false},
		"opinion only":       {"便利そう", false},
		"title with opinion": {"タイトル 便利そう", false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := isShareOnly(tc.comment, "タイトル"); got != tc.want {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}
