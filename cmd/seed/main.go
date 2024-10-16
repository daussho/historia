package main

import (
	"log"
	"time"

	"github.com/daussho/historia/domain/user"
	"github.com/daussho/historia/internal/db"
	"github.com/daussho/historia/utils"
	"github.com/daussho/historia/utils/password"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	gormDB := db.Init()

	pwd, _ := password.Hash("admin")
	admin := user.User{
		ID:        uuid.NewString(),
		Name:      "admin",
		Email:     "admin@admin.com",
		Password:  pwd,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	gormDB.Create(&admin)
	log.Println(utils.JsonStringify(admin))

	userToken := user.UserToken{
		UserID:    uuid.NewString(),
		Token:     uuid.NewString(),
		ExpiredAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	gormDB.Create(&userToken)
	log.Println(utils.JsonStringify(userToken))
}
