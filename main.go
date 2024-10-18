package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/daussho/historia/domain/healthcheck"
	"github.com/daussho/historia/domain/history"
	"github.com/daussho/historia/domain/user"
	"github.com/daussho/historia/internal/db"
	"github.com/daussho/historia/internal/middleware"
	"github.com/daussho/historia/internal/trace"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/redirect"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	log.Println("Starting application...")

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, err: ", err.Error())
	}

	sqlDB := db.Init()

	healthcheckHandler := healthcheck.NewHandler(sqlDB)

	userRepo := user.NewRepository(sqlDB)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	historyRepo := history.NewRepository(sqlDB)
	historySvc := history.NewService(historyRepo)
	historyHandler := history.NewHandler(historySvc, userService)

	mw := middleware.NewMiddleware(userService)

	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})
	app.Use(
		mw.PanicRecovery(),
		// mw.RateLimit(),
		cors.New(cors.Config{AllowOrigins: "*"}),
	)

	app.Static("/public", "./public")

	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/": "/history",
		},
		StatusCode: 301,
	}))

	app.Get("/healthcheck", trace.FiberHandler(healthcheckHandler.Healthcheck))

	app.Get("/login", trace.FiberHandler(userHandler.Login))
	app.Post("/login", trace.FiberHandler(userHandler.Login))

	app.Use("/history", mw.AuthWeb())
	app.Get("/history", trace.FiberHandler(historyHandler.ListHistory))

	apiRoute := app.Group("/api").Use(mw.AuthApi())

	apiRoute.Post("/history", trace.FiberHandler(historyHandler.SaveVisit))
	apiRoute.Put("/history/:id", trace.FiberHandler(historyHandler.UpdateVisit))

	go app.Listen(":3000")

	<-ctx.Done()

	// stop the server
	log.Println("Shutting down server gracefully...")
	err = app.Shutdown()
	if err != nil {
		log.Fatalln("failed to shutdown server: ", err.Error())
	}

	log.Println("Closing database...")
	err = sqlDB.Close()
	if err != nil {
		log.Println("failed to close database: ", err.Error())
	}

	log.Println("Application stopped")
	os.Exit(0)
}
