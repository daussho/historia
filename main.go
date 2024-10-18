package main

import (
	"log"

	"github.com/daussho/historia/domain/healthcheck"
	"github.com/daussho/historia/domain/history"
	"github.com/daussho/historia/domain/user"
	"github.com/daussho/historia/internal/db"
	"github.com/daussho/historia/internal/logger"
	"github.com/daussho/historia/internal/middleware"
	"github.com/daussho/historia/internal/trace"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Info().Msgf("Error loading .env file, err: %v", err)
	}

	gormDB := db.InitGorm()
	sqlDB := db.Init()

	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})
	app.Use(
		cors.New(cors.Config{AllowOrigins: "*"}),
		middleware.RateLimit(),
	)

	app.Static("/public", "./public")

	healthcheckHandler := healthcheck.NewHandler(gormDB)

	userService := user.NewService(gormDB)
	userHandler := user.NewHandler(userService)

	historyRepo := history.NewRepository(sqlDB)
	historySvc := history.NewService(historyRepo)
	historyHandler := history.NewHandler(historySvc, userService)

	app.Get("/healthcheck", trace.FiberHandler(healthcheckHandler.Healthcheck))

	app.Get("/login", trace.FiberHandler(userHandler.Login))
	app.Post("/login", trace.FiberHandler(userHandler.Login))

	app.Use("/history", middleware.AuthWeb(gormDB))
	app.Get("/history", trace.FiberHandler(historyHandler.ListHistory))

	apiRoute := app.Group("/api").
		Use(middleware.AuthApi(gormDB))

	apiRoute.Post("/history", trace.FiberHandler(historyHandler.SaveVisit))
	apiRoute.Put("/history/:id", trace.FiberHandler(historyHandler.UpdateVisit))

	log.Fatal(app.Listen(":3000"))
}
