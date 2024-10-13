package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/daussho/historia/domain/history"
	"github.com/daussho/historia/domain/user"
	"github.com/daussho/historia/utils/password"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init() *gorm.DB {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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
	}

	// create user token
	db.Model(&user.UserToken{}).Count(&count)
	if count == 0 {
		db.Create(&user.UserToken{
			UserID:    userID,
			Token:     uuid.NewString(),
			ExpiredAt: time.Now().AddDate(1, 0, 0).Unix(),
		})
	}
}
