package database

import (
	"log"

	"golang_world_cup/internal/user"
	"golang_world_cup/internal/match"
	"golang_world_cup/internal/bet"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// Auto Migrate all domain entities
	err = DB.AutoMigrate(
		&user.User{},
		&match.Match{},
		&bet.Bet{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connection established and migrated.")
	return DB
}
