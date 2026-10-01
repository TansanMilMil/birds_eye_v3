package ai

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func newTestJevServer(t *testing.T, status int, body string, captured *jevRequest) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("unexpected Authorization header: %q", got)
		}
		if captured != nil {
			if err := json.NewDecoder(r.Body).Decode(captured); err != nil {
				t.Errorf("failed to decode request: %v", err)
			}
		}
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
}

func newTestJevClient(baseURL string) jevClient {
	return jevClient{apiKey: "test-key", baseURL: baseURL, model: JEV_MODEL}
}

func articleAnswers(category string, confidence, release, broad, action, scope float64) string {
	body, _ := json.Marshal(jevResponse{Answers: map[string]jevAnswer{
		questionCategory:       {Choice: category, Confidence: confidence},
		questionNewRelease:     {Noul: release},
		questionBroadDevImpact: {Noul: broad},
		questionActionRequired: {Noul: action},
		questionImpactScope:    {Score: scope},
	}})
	return string(body)
}

func TestNewJevClient(t *testing.T) {
	original := os.Getenv("BIRDSEYE_TYPESAFE_API_KEY")
	defer os.Setenv("BIRDSEYE_TYPESAFE_API_KEY", original)

	os.Setenv("BIRDSEYE_TYPESAFE_API_KEY", "env-key")
	c := newJevClient()

	if c.apiKey != "env-key" {
		t.Errorf("Expected apiKey to be 'env-key', got '%s'", c.apiKey)
	}
	if c.baseURL != JEV_ENDPOINT {
		t.Errorf("Expected baseURL to be '%s', got '%s'", JEV_ENDPOINT, c.baseURL)
	}
	if c.model != JEV_MODEL {
		t.Errorf("Expected model to be '%s', got '%s'", JEV_MODEL, c.model)
	}
}

func TestJevArticleAnalyzer_Analyze_NoAPIKey(t *testing.T) {
	a := &JevArticleAnalyzer{client: jevClient{baseURL: JEV_ENDPOINT, model: JEV_MODEL}}

	if _, err := a.Analyze("title", "summary", "Zenn"); err == nil {
		t.Error("Expected error when API key is not set")
	}
}

func TestJevArticleAnalyzer_Analyze_Success(t *testing.T) {
	var req jevRequest
	server := newTestJevServer(t, http.StatusOK, articleAnswers("security_incident", 0.9, 1, 1, 1, 3), &req)
	defer server.Close()

	result, err := (&JevArticleAnalyzer{client: newTestJevClient(server.URL)}).Analyze("漏えい", "要約", "ZDNet Japan")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Category != CategorySecurity {
		t.Errorf("Expected category %q, got %q", CategorySecurity, result.Category)
	}
	if result.CategoryConfidence != 0.9 {
		t.Errorf("Expected confidence 0.9, got %v", result.CategoryConfidence)
	}
	if math.Abs(result.Importance-1) > 1e-9 {
		t.Errorf("Expected importance 1 when every signal is maximal, got %v", result.Importance)
	}

	if req.Model != JEV_MODEL {
		t.Errorf("Expected model %q, got %q", JEV_MODEL, req.Model)
	}
	state, _ := req.State.(map[string]any)
	if state["title"] != "漏えい" || state["summary"] != "要約" || state["source"] != "ZDNet Japan" {
		t.Errorf("unexpected state: %+v", req.State)
	}
	expectedTypes := map[string]string{
		questionCategory:       "choice",
		questionNewRelease:     "noul",
		questionBroadDevImpact: "noul",
		questionActionRequired: "noul",
		questionImpactScope:    "score",
	}
	if len(req.Questions) != len(expectedTypes) {
		t.Errorf("Expected %d questions in one request, got %d", len(expectedTypes), len(req.Questions))
	}
	for key, typ := range expectedTypes {
		if req.Questions[key].Type != typ {
			t.Errorf("Expected question %q to be %q, got %+v", key, typ, req.Questions[key])
		}
	}
	if criteria, _ := req.Questions[questionCategory].Criteria.(map[string]any); len(criteria) != len(categoryDescriptions) {
		t.Errorf("Expected %d category criteria, got %d", len(categoryDescriptions), len(criteria))
	}
	if criteria, _ := req.Questions[questionImpactScope].Criteria.([]any); len(criteria) != len(impactScopeLevels) {
		t.Errorf("Expected %d impact scope levels, got %d", len(impactScopeLevels), len(criteria))
	}
}

func TestJevArticleAnalyzer_Analyze_Importance(t *testing.T) {
	cases := map[string]struct {
		release, broad, action, scope float64
		want                          float64
	}{
		"all minimal":         {0, 0, 0, 0, 0},
		"weighted mix":        {1, 0.5, 0, 1.5, 0.2 + 0.15 + 0.15},
		"out of range values": {1.2, -0.1, 0, 5, 0.2 + 0.3},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			server := newTestJevServer(t, http.StatusOK,
				articleAnswers("software_dev", 0.9, tc.release, tc.broad, tc.action, tc.scope), nil)
			defer server.Close()

			result, err := (&JevArticleAnalyzer{client: newTestJevClient(server.URL)}).Analyze("t", "s", "Zenn")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(result.Importance-tc.want) > 1e-9 {
				t.Errorf("Expected importance %v, got %v", tc.want, result.Importance)
			}
		})
	}
}

func TestJevArticleAnalyzer_Analyze_LowConfidenceFallsBackToOther(t *testing.T) {
	server := newTestJevServer(t, http.StatusOK, articleAnswers("software_dev", 0.1, 0, 0, 0, 0), nil)
	defer server.Close()

	result, err := (&JevArticleAnalyzer{client: newTestJevClient(server.URL)}).Analyze("t", "s", "Zenn")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Category != CategoryOther {
		t.Errorf("Expected category %q, got %q", CategoryOther, result.Category)
	}
	if result.CategoryConfidence != 0.1 {
		t.Errorf("Expected raw confidence 0.1 to be kept, got %v", result.CategoryConfidence)
	}
}

func TestJevArticleAnalyzer_Analyze_Errors(t *testing.T) {
	cases := map[string]struct {
		status int
		body   string
	}{
		"API error status":    {http.StatusUnauthorized, `{"error":"unauthorized"}`},
		"missing all answers": {http.StatusOK, `{"answers":{}}`},
		"missing one answer":  {http.StatusOK, `{"answers":{"category":{"choice":"software_dev","confidence":0.9}}}`},
		"unknown category":    {http.StatusOK, articleAnswers("nope", 0.9, 0, 0, 0, 0)},
		"malformed payload":   {http.StatusOK, `not json`},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			server := newTestJevServer(t, tc.status, tc.body, nil)
			defer server.Close()

			if _, err := (&JevArticleAnalyzer{client: newTestJevClient(server.URL)}).Analyze("t", "s", "Zenn"); err == nil {
				t.Error("Expected error")
			}
		})
	}
}
