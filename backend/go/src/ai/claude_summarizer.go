package ai

import (
	"fmt"
	"os"
)

const CLAUDE_CHAT_ENDPOINT = "https://api.anthropic.com/v1/messages"
const CLAUDE_MODEL = "claude-3-5-sonnet-20241022"

type ClaudeSummarizer struct {
	apiKey      string
	baseURL     string
	claudeModel string
	maxTokens   int
}

type ClaudeRequest struct {
	Model     string          `json:"model"`
	Messages  []ClaudeMessage `json:"messages"`
	MaxTokens int             `json:"max_tokens"`
}

type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func NewClaudeSummarizer() *ClaudeSummarizer {
	return &ClaudeSummarizer{
		apiKey:      os.Getenv("BIRDSEYE_BIRDSEYEAPI_V2_CLAUDE_API_KEY"),
		baseURL:     CLAUDE_CHAT_ENDPOINT,
		claudeModel: CLAUDE_MODEL,
		maxTokens:   1024,
	}
}

func (s *ClaudeSummarizer) Summarize(text string) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("claude API key not found")
	}

	reqBody := ClaudeRequest{
		Model:     s.claudeModel,
		Messages:  []ClaudeMessage{{Role: "user", Content: summarizePrompt(text)}},
		MaxTokens: s.maxTokens,
	}
	headers := map[string]string{
		"x-api-key":         s.apiKey,
		"anthropic-version": "2023-06-01",
	}

	var respBody ClaudeResponse
	if err := postJSON(s.baseURL, headers, reqBody, &respBody); err != nil {
		return "", err
	}

	if len(respBody.Content) == 0 {
		return "", fmt.Errorf("no summary was generated")
	}

	return respBody.Content[0].Text, nil
}
