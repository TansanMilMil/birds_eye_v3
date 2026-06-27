package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(DBFilePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(GetMigrationModels()...)
	if err != nil {
		return nil, err
	}

	return db, nil
}
