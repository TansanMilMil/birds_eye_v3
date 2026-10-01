package ai

import (
	"encoding/json"
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

func newTestJevCategorizer(baseURL string) *JevCategorizer {
	return &JevCategorizer{apiKey: "test-key", baseURL: baseURL, model: JEV_MODEL}
}

func TestNewJevCategorizer(t *testing.T) {
	original := os.Getenv("BIRDSEYE_TYPESAFE_API_KEY")
	defer os.Setenv("BIRDSEYE_TYPESAFE_API_KEY", original)

	os.Setenv("BIRDSEYE_TYPESAFE_API_KEY", "env-key")
	c := NewJevCategorizer()

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

func TestJevCategorizer_Categorize_NoAPIKey(t *testing.T) {
	c := &JevCategorizer{baseURL: JEV_ENDPOINT, model: JEV_MODEL}

	if _, err := c.Categorize("title", "summary", "Zenn"); err == nil {
		t.Error("Expected error when API key is not set")
	}
}

func TestJevCategorizer_Categorize_Success(t *testing.T) {
	var req jevRequest
	server := newTestJevServer(t, http.StatusOK,
		`{"answers":{"category":{"choice":"security_incident","confidence":0.9}}}`, &req)
	defer server.Close()

	result, err := newTestJevCategorizer(server.URL).Categorize("漏えい", "要約", "ZDNet Japan")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Category != CategorySecurity {
		t.Errorf("Expected category %q, got %q", CategorySecurity, result.Category)
	}
	if result.Confidence != 0.9 {
		t.Errorf("Expected confidence 0.9, got %v", result.Confidence)
	}

	if req.Model != JEV_MODEL {
		t.Errorf("Expected model %q, got %q", JEV_MODEL, req.Model)
	}
	if req.State.Title != "漏えい" || req.State.Summary != "要約" || req.State.Source != "ZDNet Japan" {
		t.Errorf("unexpected state: %+v", req.State)
	}
	q, ok := req.Questions["category"]
	if !ok || q.Type != "choice" {
		t.Fatalf("expected a choice question named category, got %+v", req.Questions)
	}
	if len(q.Criteria) != len(categoryDescriptions) {
		t.Errorf("Expected %d criteria, got %d", len(categoryDescriptions), len(q.Criteria))
	}
}

func TestJevCategorizer_Categorize_LowConfidenceFallsBackToOther(t *testing.T) {
	server := newTestJevServer(t, http.StatusOK,
		`{"answers":{"category":{"choice":"software_dev","confidence":0.1}}}`, nil)
	defer server.Close()

	result, err := newTestJevCategorizer(server.URL).Categorize("t", "s", "Zenn")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Category != CategoryOther {
		t.Errorf("Expected category %q, got %q", CategoryOther, result.Category)
	}
	if result.Confidence != 0.1 {
		t.Errorf("Expected raw confidence 0.1 to be kept, got %v", result.Confidence)
	}
}

func TestJevCategorizer_Categorize_Errors(t *testing.T) {
	cases := map[string]struct {
		status int
		body   string
	}{
		"API error status":  {http.StatusUnauthorized, `{"error":"unauthorized"}`},
		"missing answer":    {http.StatusOK, `{"answers":{}}`},
		"unknown category":  {http.StatusOK, `{"answers":{"category":{"choice":"nope","confidence":0.9}}}`},
		"malformed payload": {http.StatusOK, `not json`},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			server := newTestJevServer(t, tc.status, tc.body, nil)
			defer server.Close()

			if _, err := newTestJevCategorizer(server.URL).Categorize("t", "s", "Zenn"); err == nil {
				t.Error("Expected error")
			}
		})
	}
}
