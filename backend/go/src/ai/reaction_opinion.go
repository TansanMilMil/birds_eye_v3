package ai

import "fmt"

const OPINION_COMMENTS_PER_REQUEST = 20

// Comments are dropped only when Jev is fairly sure they carry no opinion;
// uncertain ones are kept so that real reactions are not lost.
const MIN_OPINION_PROBABILITY = 0.4

type ReactionOpinionClassifier interface {
	HasOpinion(articleTitle string, comments []string) ([]bool, error)
}

type JevReactionOpinionClassifier struct {
	client jevClient
}

func NewJevReactionOpinionClassifier() *JevReactionOpinionClassifier {
	return &JevReactionOpinionClassifier{client: newJevClient()}
}

func (c *JevReactionOpinionClassifier) HasOpinion(articleTitle string, comments []string) ([]bool, error) {
	result := make([]bool, len(comments))

	for start := 0; start < len(comments); start += OPINION_COMMENTS_PER_REQUEST {
		end := min(start+OPINION_COMMENTS_PER_REQUEST, len(comments))
		state := reactionState{ArticleTitle: articleTitle, Comments: map[string]string{}}
		questions := map[string]jevQuestion{}

		for i, comment := range comments[start:end] {
			key := commentKey(i)
			state.Comments[key] = comment
			questions[key] = jevQuestion{
				Type: "noul",
				Instructions: fmt.Sprintf(
					"Does the comment %q in comments express the writer's own opinion, impression, evaluation, or question about the article or its subject? "+
						"Answer no if it only repeats the article title, shares a link, quotes the article, or is a plain retweet without any words of the writer's own.",
					key,
				),
			}
		}

		answers, err := c.client.ask(state, questions)
		if err != nil {
			return nil, err
		}

		for i := range comments[start:end] {
			result[start+i] = answers[commentKey(i)].Noul >= MIN_OPINION_PROBABILITY
		}
	}

	return result, nil
}
