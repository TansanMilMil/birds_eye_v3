package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	newsHandler := NewNewsHandler(db)
	trendHandler := NewTrendHandler()
	feedHandler := NewFeedHandler(db)

	newsRoutes := r.Group("/news")
	{
		newsRoutes.GET("/:target_date", newsHandler.GetNewsByDate)
		newsRoutes.GET("/news-reactions/:news-id", newsHandler.GetNewsReactionsById)
		newsRoutes.POST("/scrape", newsHandler.Scrape)

		newsRoutes.GET("/trends", trendHandler.GetTrends)
		newsRoutes.GET("/rss", feedHandler.GetRSS)
	}
}
