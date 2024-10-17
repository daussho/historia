package main

import (
	"log"

	"github.com/daussho/historia/domain/healthcheck"
	"github.com/daussho/historia/domain/history"
	"github.com/daussho/historia/internal/db"
	"github.com/daussho/historia/internal/middleware"
	"github.com/daussho/historia/internal/trace"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	gormDB := db.Init()

	healthcheckHandler := healthcheck.NewHandler(gormDB)
	historySvc := history.NewService(gormDB)
	historyHandler := history.NewHandler(historySvc)

	app := fiber.New()
	app.Use(
		cors.New(cors.Config{AllowOrigins: "*"}),
		middleware.RateLimit(),
	)

	app.Static("/public", "./public")

	app.Get("/healthcheck", trace.FiberHandler(healthcheckHandler.Healthcheck))

	apiRoute := app.Group("/api").
		Use(middleware.AuthApi(gormDB))

	apiRoute.Post("/history", trace.FiberHandler(historyHandler.SaveVisit))
	apiRoute.Put("/history/:id", trace.FiberHandler(historyHandler.UpdateVisit))

	log.Fatal(app.Listen(":3000"))
}
