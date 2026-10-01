package ai

import (
	"fmt"
)

const MIN_CATEGORY_CONFIDENCE = 0.4

const categoryInstructions = "Choose the single category that best matches the main subject of the article. " +
	"Pick AI categories only when AI itself is the main subject, not when it is merely a means. " +
	"If an article is mainly about a security incident or security measure, choose the security category even if AI is involved. " +
	"Press releases from companies announcing a product or service belong to the product release category."

const (
	questionCategory       = "category"
	questionNewRelease     = "new_release"
	questionBroadDevImpact = "broad_dev_impact"
	questionActionRequired = "action_required"
	questionImpactScope    = "impact_scope"
)

var impactScopeLevels = []string{
	"Niche: concerns an individual, a personal note, or a single small project",
	"Limited: concerns the users of a specific product or a specific community",
	"Broad: concerns a large part of software developers or the IT industry",
	"Industry-wide: major news that affects the whole IT industry or society",
}

var importanceWeights = map[string]float64{
	questionNewRelease:     0.2,
	questionBroadDevImpact: 0.3,
	questionActionRequired: 0.2,
	questionImpactScope:    0.3,
}

type articleState struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Source  string `json:"source"`
}

type JevArticleAnalyzer struct {
	client jevClient
}

func NewJevArticleAnalyzer() *JevArticleAnalyzer {
	return &JevArticleAnalyzer{client: newJevClient()}
}

func (a *JevArticleAnalyzer) Analyze(title, summary, sourceBy string) (ArticleAnalysis, error) {
	answers, err := a.client.ask(
		articleState{Title: title, Summary: summary, Source: sourceBy},
		articleQuestions(),
	)
	if err != nil {
		return ArticleAnalysis{}, err
	}

	category := answers[questionCategory]
	if _, known := categoryDescriptions[category.Choice]; !known {
		return ArticleAnalysis{}, fmt.Errorf("unknown category returned: %q", category.Choice)
	}

	result := ArticleAnalysis{
		Category:           category.Choice,
		CategoryConfidence: category.Confidence,
		Importance:         combineImportance(answers),
	}
	if category.Confidence < MIN_CATEGORY_CONFIDENCE {
		result.Category = CategoryOther
	}
	return result, nil
}

func articleQuestions() map[string]jevQuestion {
	return map[string]jevQuestion{
		questionCategory: {
			Type:         "choice",
			Instructions: categoryInstructions,
			Criteria:     categoryDescriptions,
		},
		questionNewRelease: {
			Type:         "noul",
			Instructions: "Does the article announce a new release, launch, or official announcement of a product, service, model, or version?",
		},
		questionBroadDevImpact: {
			Type:         "noul",
			Instructions: "Is the news likely to affect the work of many software developers or IT engineers?",
		},
		questionActionRequired: {
			Type:         "noul",
			Instructions: "Does the article report a breaking change, a vulnerability, or a security incident that readers may need to act on?",
		},
		questionImpactScope: {
			Type:         "score",
			Instructions: "How wide is the scope of people and organizations affected by the news in this article?",
			Criteria:     impactScopeLevels,
		},
	}
}

// combineImportance normalizes every signal to 0..1 and returns their weighted sum.
func combineImportance(answers map[string]jevAnswer) float64 {
	signals := map[string]float64{
		questionNewRelease:     answers[questionNewRelease].Noul,
		questionBroadDevImpact: answers[questionBroadDevImpact].Noul,
		questionActionRequired: answers[questionActionRequired].Noul,
		questionImpactScope:    answers[questionImpactScope].Score / float64(len(impactScopeLevels)-1),
	}

	importance := 0.0
	for key, weight := range importanceWeights {
		importance += weight * clamp01(signals[key])
	}
	return importance
}

func clamp01(v float64) float64 {
	return min(max(v, 0), 1)
}
