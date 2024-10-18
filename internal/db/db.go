package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGorm() *gorm.DB {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, database)
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db, err := gormDB.DB()
	if err != nil {
		log.Fatalf("failed to get database: %v", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(2)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(time.Hour)

	return gormDB
}

func Init() *sqlx.DB {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", user, password, host, port, database))
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	return db
}
