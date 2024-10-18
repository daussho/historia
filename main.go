package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	logger.Log().Info("Starting application...")

	err := godotenv.Load()
	if err != nil {
		logger.Log().Infof("Error loading .env file, err: %v", err)
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

	healthcheckHandler := healthcheck.NewHandler(sqlDB)

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

	apiRoute := app.Group("/api").Use(middleware.AuthApi(gormDB))

	apiRoute.Post("/history", trace.FiberHandler(historyHandler.SaveVisit))
	apiRoute.Put("/history/:id", trace.FiberHandler(historyHandler.UpdateVisit))

	go app.Listen(":3000")

	<-ctx.Done()

	// stop the server
	logger.Log().Info("Shutting down server gracefully...")
	err = app.Shutdown()
	if err != nil {
		logger.Log().Fatalf("failed to shutdown server: %v", err)
	}

	logger.Log().Info("Closing database...")
	err = sqlDB.Close()
	if err != nil {
		logger.Log().Errorf("failed to close database: %v", err)
	}

	logger.Log().Info("Application stopped")
	os.Exit(0)
}
