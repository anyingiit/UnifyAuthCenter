package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Sessions struct {
	UUID      uuid.UUID `gorm:"primarykey;type:string"`
	CreatedAt time.Time
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}
	db.AutoMigrate(&Sessions{})
	db.Create(&Sessions{
		UUID:      uuid.New(),
		CreatedAt: time.Time{},
	})

	var sessions []Sessions

	result := db.Find(&sessions)
	log.Println(result.RowsAffected)
	log.Println(result.Error)

	for _, session := range sessions {
		log.Println(session.UUID, session.CreatedAt)
	}
}
