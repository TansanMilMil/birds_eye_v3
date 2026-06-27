package trend

import "github.com/birdseyeapi/birds_eye_v3/go/src/models"

type TrendCatcher interface {
	GetTrends() ([]models.News, error)
}
