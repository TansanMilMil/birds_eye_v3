package db

import "github.com/birdseyeapi/birds_eye_v3/go/src/models"

func GetMigrationModels() []interface{} {
	return []interface{}{
		&models.News{},
		&models.NewsReaction{},
	}
}
