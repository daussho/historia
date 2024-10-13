package main

import (
	"log"

	"github.com/daussho/historia/domain/history"
	"github.com/daussho/historia/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	gormDB := db.Init()
	db.Migrate(gormDB)
	db.Seed(gormDB)

	historySvc := history.NewService(gormDB)
	historyHandler := history.NewHandler(historySvc)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Static("/public", "./public")

	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return gormDB.Raw("SELECT 1").Error
	})

	apiRoute := app.Group("/api")
	apiRoute.Post("/history", historyHandler.SaveVisit)
	apiRoute.Put("/history/:id", historyHandler.UpdateVisit)

	log.Fatal(app.Listen(":3000"))
}
