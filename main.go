package main

import (
	"log"

	"github.com/daussho/historia/domain/history"
	"github.com/daussho/historia/internal/db"
	"github.com/gofiber/fiber/v2"
)

func main() {
	gormDB := db.Init()
	db.Migrate(gormDB)

	historySvc := history.NewService(gormDB)
	historyHandler := history.NewHandler(historySvc)

	app := fiber.New()

	apiRoute := app.Group("/api")
	apiRoute.Post("/history", historyHandler.SaveVisit)

	log.Fatal(app.Listen(":3000"))
}
