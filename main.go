package main

import (
	"log"

	"github.com/daussho/historia/domain/history"
	"github.com/daussho/historia/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
		err := gormDB.Raw("SELECT 1").Error
		if err != nil {
			return c.SendStatus(500)
		}

		c.SendString("OK")
		return c.SendStatus(200)
	})

	apiRoute := app.Group("/api")
	apiRoute.Post("/history", historyHandler.SaveVisit)
	apiRoute.Put("/history/:id", historyHandler.UpdateVisit)

	log.Fatal(app.Listen(":3000"))
}
