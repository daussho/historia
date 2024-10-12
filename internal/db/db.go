package db

import (
	"log"

	"github.com/daussho/historia/domain/history"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("historia.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&history.History{})
}
