package models

import (
	"time"

	"gorm.io/gorm"
)

type ReactionSentiment struct {
	Positive int `json:"positive"`
	Neutral  int `json:"neutral"`
	Negative int `json:"negative"`
}

type NewsReaction struct {
	gorm.Model
	ID              uint      `gorm:"primarykey"`
	NewsID          uint      `gorm:"index" json:"newsId"`
	Author          string    `json:"author"`
	Comment         string    `json:"comment"`
	ScrapedDateTime time.Time `json:"scrapedDateTime"`
	CommentUrl      string    `json:"commentUrl"`
}
