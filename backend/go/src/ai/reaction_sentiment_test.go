package ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newSentimentServer answers each comment with the score registered for its text.
func newSentimentServer(t *testing.T, scores map[string]float64, requestSizes *[]int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			State     reactionState          `json:"state"`
			Questions map[string]jevQuestion `json:"questions"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}
		*requestSizes = append(*requestSizes, len(req.Questions))

		answers := map[string]jevAnswer{}
		for key, q := range req.Questions {
			if q.Type != "score" {
				t.Errorf("Expected score question, got %q", q.Type)
			}
			comment, ok := req.State.Comments[key]
			if !ok {
				t.Errorf("question %q has no matching comment in state", key)
			}
			answers[key] = jevAnswer{Score: scores[comment]}
		}
		json.NewEncoder(w).Encode(jevResponse{Answers: answers})
	}))
}

func TestJevReactionSentimentAnalyzer_Analyze(t *testing.T) {
	scores := map[string]float64{"bad": 0.2, "meh": 1.4, "good": 1.6, "low": -0.3, "high": 2.7}
	var sizes []int
	server := newSentimentServer(t, scores, &sizes)
	defer server.Close()

	a := &JevReactionSentimentAnalyzer{client: newTestJevClient(server.URL)}
	counts, err := a.Analyze("title", []string{"bad", "meh", "good", "low", "high"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := SentimentCounts{Positive: 2, Neutral: 1, Negative: 2}
	if counts != want {
		t.Errorf("Expected %+v, got %+v", want, counts)
	}
}

func TestJevReactionSentimentAnalyzer_Analyze_SplitsRequests(t *testing.T) {
	comments := make([]string, SENTIMENT_COMMENTS_PER_REQUEST+5)
	scores := map[string]float64{}
	for i := range comments {
		comments[i] = fmt.Sprintf("comment-%d", i)
		scores[comments[i]] = 2
	}
	var sizes []int
	server := newSentimentServer(t, scores, &sizes)
	defer server.Close()

	a := &JevReactionSentimentAnalyzer{client: newTestJevClient(server.URL)}
	counts, err := a.Analyze("title", comments)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sizes) != 2 || sizes[0] != SENTIMENT_COMMENTS_PER_REQUEST || sizes[1] != 5 {
		t.Errorf("unexpected request sizes: %v", sizes)
	}
	if counts.Positive != len(comments) {
		t.Errorf("Expected all %d comments positive, got %+v", len(comments), counts)
	}
}

func TestJevReactionSentimentAnalyzer_Analyze_NoComments(t *testing.T) {
	a := &JevReactionSentimentAnalyzer{client: jevClient{}}

	counts, err := a.Analyze("title", nil)
	if err != nil || counts != (SentimentCounts{}) {
		t.Errorf("Expected empty counts without request, got %+v, %v", counts, err)
	}
}

func TestJevReactionSentimentAnalyzer_Analyze_Error(t *testing.T) {
	server := newTestJevServer(t, http.StatusOK, `{"answers":{}}`, nil)
	defer server.Close()

	a := &JevReactionSentimentAnalyzer{client: newTestJevClient(server.URL)}
	if _, err := a.Analyze("title", []string{"a"}); err == nil {
		t.Error("Expected error when an answer is missing")
	}
}
