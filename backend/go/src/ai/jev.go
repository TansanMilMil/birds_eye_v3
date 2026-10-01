package ai

import (
	"fmt"
	"os"
)

const JEV_ENDPOINT = "https://api.typesafe.ai/v1/systemone"
const JEV_MODEL = "jev-latest"

type jevClient struct {
	apiKey  string
	baseURL string
	model   string
}

type jevRequest struct {
	Model     string                 `json:"model"`
	State     any                    `json:"state"`
	Questions map[string]jevQuestion `json:"questions"`
}

type jevQuestion struct {
	Type         string `json:"type"`
	Instructions string `json:"instructions"`
	Criteria     any    `json:"criteria,omitempty"`
}

type jevAnswer struct {
	Choice     string  `json:"choice"`
	Score      float64 `json:"score"`
	Noul       float64 `json:"noul"`
	Confidence float64 `json:"confidence"`
}

type jevResponse struct {
	Answers map[string]jevAnswer `json:"answers"`
}

func newJevClient() jevClient {
	return jevClient{
		apiKey:  os.Getenv("BIRDSEYE_TYPESAFE_API_KEY"),
		baseURL: JEV_ENDPOINT,
		model:   JEV_MODEL,
	}
}

func (c jevClient) ask(state any, questions map[string]jevQuestion) (map[string]jevAnswer, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("typesafe API key not found")
	}

	reqBody := jevRequest{Model: c.model, State: state, Questions: questions}
	headers := map[string]string{"Authorization": "Bearer " + c.apiKey}

	var respBody jevResponse
	if err := postJSON(c.baseURL, headers, reqBody, &respBody); err != nil {
		return nil, err
	}

	for key := range questions {
		if _, ok := respBody.Answers[key]; !ok {
			return nil, fmt.Errorf("no answer was returned for %q", key)
		}
	}
	return respBody.Answers, nil
}
