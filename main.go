package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/daussho/historia/domain/common"
	"github.com/daussho/historia/domain/healthcheck"
	"github.com/daussho/historia/domain/history"
	"github.com/daussho/historia/domain/user"
	"github.com/daussho/historia/internal/db"
	"github.com/daussho/historia/internal/trace"
	context_util "github.com/daussho/historia/utils/context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
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
	)
	app.Use(authMiddleware(gormDB))

	app.Static("/public", "./public")

	app.Get("/healthcheck", trace.FiberHandler(healthcheckHandler.Healthcheck))

	apiRoute := app.Group("/api")
	apiRoute.Post("/history", historyHandler.SaveVisit)
	apiRoute.Put("/history/:id", historyHandler.UpdateVisit)

	log.Fatal(app.Listen(":3000"))
}

type fiberHandler = func(ctx *fiber.Ctx) error

func authMiddleware(db *gorm.DB) fiberHandler {
	return func(ctx *fiber.Ctx) error {
		path := ctx.Request().URI().Path()

		if string(path) == "/healthcheck" {
			return ctx.Next()
		}

		if strings.Contains(string(path), "/public") {
			return ctx.Next()
		}

		headers := ctx.GetReqHeaders()

		tokens := headers["Authorization"]
		if len(tokens) == 0 {
			return common.NewErrorResponse(ctx, fiber.StatusUnauthorized, fmt.Errorf("token empty"))
		}

		bearerToken := tokens[0]
		if bearerToken == "" {
			return common.NewErrorResponse(ctx, fiber.StatusUnauthorized, fmt.Errorf("token empty"))
		}

		tokens = strings.Split(bearerToken, " ")
		if len(tokens) != 2 || tokens[0] != "Bearer" {
			return common.NewErrorResponse(ctx, fiber.StatusUnauthorized, fmt.Errorf("invalid token"))
		}

		token := tokens[1]
		if token == "" {
			return common.NewErrorResponse(ctx, fiber.StatusUnauthorized, fmt.Errorf("token empty"))
		}

		var userToken user.UserToken
		err := db.WithContext(ctx.Context()).First(&userToken, "token = ?", token).Error
		if err != nil {
			return common.NewErrorResponse(ctx, fiber.StatusUnauthorized, fmt.Errorf("token not found"))
		}

		var user user.User
		err = db.WithContext(ctx.Context()).First(&user, "id = ?", userToken.UserID).Error
		if err != nil {
			return common.NewErrorResponse(ctx, fiber.StatusUnauthorized, fmt.Errorf("user id %s not found", userToken.UserID))
		}

		context_util.SetUserCtx(ctx, user)

		return ctx.Next()
	}
}
