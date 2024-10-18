package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/daussho/historia/internal/logger"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Init() *sqlx.DB {
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		DBName:               os.Getenv("DB_DATABASE"),
		Params:               map[string]string{},
		Logger:               nil,
		CheckConnLiveness:    true,
		ParseTime:            true,
		Collation:            "utf8mb4_general_ci",
		AllowNativePasswords: true,
	}

	dsn := cfg.FormatDSN()
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	logger.Log().Info("sqlx connected to database")

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(time.Minute)
	db.SetConnMaxLifetime(time.Minute)

	return db
}
