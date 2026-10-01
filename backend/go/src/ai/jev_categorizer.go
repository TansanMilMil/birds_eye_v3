package ai

import (
	"fmt"
	"os"
)

const JEV_ENDPOINT = "https://api.typesafe.ai/v1/systemone"
const JEV_MODEL = "jev-latest"
const MIN_CATEGORY_CONFIDENCE = 0.4

const categoryInstructions = "Choose the single category that best matches the main subject of the article. " +
	"Pick AI categories only when AI itself is the main subject, not when it is merely a means. " +
	"If an article is mainly about a security incident or security measure, choose the security category even if AI is involved. " +
	"Press releases from companies announcing a product or service belong to the product release category."

type JevCategorizer struct {
	apiKey  string
	baseURL string
	model   string
}

type jevRequest struct {
	Model     string                 `json:"model"`
	State     jevState               `json:"state"`
	Questions map[string]jevQuestion `json:"questions"`
}

type jevState struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Source  string `json:"source"`
}

type jevQuestion struct {
	Type         string            `json:"type"`
	Instructions string            `json:"instructions"`
	Criteria     map[string]string `json:"criteria"`
}

type jevResponse struct {
	Answers map[string]struct {
		Choice     string  `json:"choice"`
		Confidence float64 `json:"confidence"`
	} `json:"answers"`
}

func NewJevCategorizer() *JevCategorizer {
	return &JevCategorizer{
		apiKey:  os.Getenv("BIRDSEYE_TYPESAFE_API_KEY"),
		baseURL: JEV_ENDPOINT,
		model:   JEV_MODEL,
	}
}

func (c *JevCategorizer) Categorize(title, summary, sourceBy string) (CategoryResult, error) {
	if c.apiKey == "" {
		return CategoryResult{}, fmt.Errorf("typesafe API key not found")
	}

	reqBody := jevRequest{
		Model: c.model,
		State: jevState{Title: title, Summary: summary, Source: sourceBy},
		Questions: map[string]jevQuestion{
			"category": {
				Type:         "choice",
				Instructions: categoryInstructions,
				Criteria:     categoryDescriptions,
			},
		},
	}
	headers := map[string]string{"Authorization": "Bearer " + c.apiKey}

	var respBody jevResponse
	if err := postJSON(c.baseURL, headers, reqBody, &respBody); err != nil {
		return CategoryResult{}, err
	}

	answer, ok := respBody.Answers["category"]
	if !ok {
		return CategoryResult{}, fmt.Errorf("no category was returned")
	}
	if _, known := categoryDescriptions[answer.Choice]; !known {
		return CategoryResult{}, fmt.Errorf("unknown category returned: %q", answer.Choice)
	}

	category := answer.Choice
	if answer.Confidence < MIN_CATEGORY_CONFIDENCE {
		category = CategoryOther
	}

	return CategoryResult{Category: category, Confidence: answer.Confidence}, nil
}
