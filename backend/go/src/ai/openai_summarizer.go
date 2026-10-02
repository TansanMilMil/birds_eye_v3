package ai

import (
	"fmt"
	"os"
)

const BIRDSEYE_OPENAI_CHAT_ENDPOINT = "https://api.openai.com/v1/chat/completions"

type OpenAISummarizer struct {
	apiKey      string
	baseURL     string
	openAIModel string
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewOpenAISummarizer() *OpenAISummarizer {
	return &OpenAISummarizer{
		apiKey:      os.Getenv("BIRDSEYE_BIRDSEYEAPI_V2_OPENAI_API_KEY"),
		baseURL:     BIRDSEYE_OPENAI_CHAT_ENDPOINT,
		openAIModel: os.Getenv("BIRDSEYE_OPENAI_MODEL"),
	}
}

func (s *OpenAISummarizer) Summarize(text string) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("OpenAI API key not found")
	}
	if s.openAIModel == "" {
		return "", fmt.Errorf("OpenAI model not found")
	}

	reqBody := OpenAIRequest{
		Model:    s.openAIModel,
		Messages: []Message{{Role: "user", Content: summarizePrompt(text)}},
	}
	headers := map[string]string{"Authorization": "Bearer " + s.apiKey}

	var respBody OpenAIResponse
	if err := postJSON(s.baseURL, headers, reqBody, &respBody); err != nil {
		return "", err
	}

	if len(respBody.Choices) == 0 {
		return "", fmt.Errorf("no summary was generated")
	}

	return respBody.Choices[0].Message.Content, nil
}
