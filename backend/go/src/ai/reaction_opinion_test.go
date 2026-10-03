package ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func newOpinionServer(t *testing.T, nouls map[string]float64, requestSizes *[]int) *httptest.Server {
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
			if q.Type != "noul" {
				t.Errorf("Expected noul question, got %q", q.Type)
			}
			comment, ok := req.State.Comments[key]
			if !ok {
				t.Errorf("question %q has no matching comment in state", key)
			}
			answers[key] = jevAnswer{Noul: nouls[comment]}
		}
		json.NewEncoder(w).Encode(jevResponse{Answers: answers})
	}))
}

func TestJevReactionOpinionClassifier_HasOpinion(t *testing.T) {
	nouls := map[string]float64{"share": 0.1, "border": MIN_OPINION_PROBABILITY, "unsure": 0.5, "opinion": 0.9}
	var sizes []int
	server := newOpinionServer(t, nouls, &sizes)
	defer server.Close()

	c := &JevReactionOpinionClassifier{client: newTestJevClient(server.URL)}
	got, err := c.HasOpinion("title", []string{"share", "border", "unsure", "opinion"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []bool{false, true, true, true}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestJevReactionOpinionClassifier_HasOpinion_SplitsRequests(t *testing.T) {
	comments := make([]string, OPINION_COMMENTS_PER_REQUEST+5)
	nouls := map[string]float64{}
	for i := range comments {
		comments[i] = fmt.Sprintf("comment-%d", i)
		if i%2 == 0 {
			nouls[comments[i]] = 1
		}
	}
	var sizes []int
	server := newOpinionServer(t, nouls, &sizes)
	defer server.Close()

	c := &JevReactionOpinionClassifier{client: newTestJevClient(server.URL)}
	got, err := c.HasOpinion("title", comments)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sizes) != 2 || sizes[0] != OPINION_COMMENTS_PER_REQUEST || sizes[1] != 5 {
		t.Errorf("unexpected request sizes: %v", sizes)
	}
	for i, ok := range got {
		if ok != (i%2 == 0) {
			t.Errorf("comment %d: expected %v, got %v", i, i%2 == 0, ok)
		}
	}
}

func TestJevReactionOpinionClassifier_HasOpinion_NoComments(t *testing.T) {
	c := &JevReactionOpinionClassifier{client: jevClient{}}

	got, err := c.HasOpinion("title", nil)
	if err != nil || len(got) != 0 {
		t.Errorf("Expected empty result without request, got %v, %v", got, err)
	}
}

func TestJevReactionOpinionClassifier_HasOpinion_Error(t *testing.T) {
	server := newTestJevServer(t, http.StatusOK, `{"answers":{}}`, nil)
	defer server.Close()

	c := &JevReactionOpinionClassifier{client: newTestJevClient(server.URL)}
	if _, err := c.HasOpinion("title", []string{"a"}); err == nil {
		t.Error("Expected error when an answer is missing")
	}
}
