package news

import (
	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
)

type ScrapingNews interface {
	ExtractNews() ([]models.News, error)
	GetSourceBy() string
}
