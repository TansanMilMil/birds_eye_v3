package repository

import (
	"net/http"
	"time"

	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const sqliteDateTimeLayout = "2006-01-02 15:04:05"

type NewsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{
		db: db,
	}
}

// GetNews returns news created on the day starting at t, in t's location.
func (n *NewsRepository) GetNews(t time.Time, c *gin.Context) []models.News {
	var news []models.News

	// datetime() normalizes stored values with any UTC offset to UTC, so the bounds must be UTC too.
	start := t.UTC().Format(sqliteDateTimeLayout)
	end := t.AddDate(0, 0, 1).UTC().Format(sqliteDateTimeLayout)

	result := n.db.
		Where("datetime(created_at) >= ? AND datetime(created_at) < ?", start, end).
		Limit(100).
		Preload("Reactions").
		Find(&news)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return nil
	}
	return news
}

func (n *NewsRepository) GetLatestNews(limit int) ([]models.News, error) {
	var news []models.News
	result := n.db.
		Order("datetime(created_at) DESC").
		Order("id DESC").
		Limit(limit).
		Find(&news)
	return news, result.Error
}
