package database

import (
	"fmt"

	// "encoding/json"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB gorm.DB
)

func InitDB() {
	db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintln("Error when trying to make connection to database"))
	}

	db.AutoMigrate(&Card{})
	db.AutoMigrate(&Deck{})
	db.AutoMigrate(&BlackJack{})

	DB = *db
}
