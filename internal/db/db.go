package db

import (
	"log"
	"time"

	"github.com/daussho/historia/domain/history"
	"github.com/daussho/historia/domain/user"
	"github.com/daussho/historia/utils/password"
	"github.com/google/uuid"
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
	db.AutoMigrate(
		&user.User{},
		&user.UserToken{},
		&history.History{},
	)
}

func Seed(db *gorm.DB) {
	var count int64
	db.Model(&user.User{}).Count(&count)

	var userID string
	if count == 0 {
		userID = uuid.NewString()
		pwd, _ := password.Hash("admin")
		db.Create(&user.User{
			ID:       userID,
			Name:     "admin",
			Email:    "admin@admin.com",
			Password: pwd,
		})
	} else if count == 1 {
		db.Model(&user.User{}).Select("id").First(&userID)
		db.Model(&user.UserToken{}).Count(&count)
		if count == 0 {
			db.Create(&user.UserToken{
				UserID:    userID,
				Token:     uuid.NewString(),
				ExpiredAt: time.Now().AddDate(1, 0, 0).Unix(),
			})
		}
	}
}
