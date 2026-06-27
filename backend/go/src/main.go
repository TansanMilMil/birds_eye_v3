package main

import (
	"log"

	api "github.com/birdseyeapi/birds_eye_v3/go/src/api"
	"github.com/birdseyeapi/birds_eye_v3/go/src/db"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	api.Init(db)
}
