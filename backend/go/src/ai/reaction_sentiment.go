package ai

import (
	"fmt"
	"math"
)

const SENTIMENT_COMMENTS_PER_REQUEST = 20

var sentimentLevels = []string{
	"Negative: criticizes, doubts, disagrees with, or mocks the article or its subject",
	"Neutral: neither positive nor negative, such as a factual remark, a memo, or a question",
	"Positive: agrees with, praises, welcomes, or is excited about the article or its subject",
}

type SentimentCounts struct {
	Positive int
	Neutral  int
	Negative int
}

type ReactionSentimentAnalyzer interface {
	Analyze(articleTitle string, comments []string) (SentimentCounts, error)
}

type reactionState struct {
	ArticleTitle string            `json:"article_title"`
	Comments     map[string]string `json:"comments"`
}

type JevReactionSentimentAnalyzer struct {
	client jevClient
}

func NewJevReactionSentimentAnalyzer() *JevReactionSentimentAnalyzer {
	return &JevReactionSentimentAnalyzer{client: newJevClient()}
}

func (a *JevReactionSentimentAnalyzer) Analyze(articleTitle string, comments []string) (SentimentCounts, error) {
	var counts SentimentCounts

	for start := 0; start < len(comments); start += SENTIMENT_COMMENTS_PER_REQUEST {
		end := min(start+SENTIMENT_COMMENTS_PER_REQUEST, len(comments))
		state := reactionState{ArticleTitle: articleTitle, Comments: map[string]string{}}
		questions := map[string]jevQuestion{}

		for i, comment := range comments[start:end] {
			key := fmt.Sprintf("c%d", i)
			state.Comments[key] = comment
			questions[key] = jevQuestion{
				Type:         "score",
				Instructions: fmt.Sprintf("What stance does the comment %q in comments take toward the article?", key),
				Criteria:     sentimentLevels,
			}
		}

		answers, err := a.client.ask(state, questions)
		if err != nil {
			return SentimentCounts{}, err
		}

		for key := range questions {
			switch min(max(math.Round(answers[key].Score), 0), 2) {
			case 0:
				counts.Negative++
			case 2:
				counts.Positive++
			default:
				counts.Neutral++
			}
		}
	}

	return counts, nil
}
