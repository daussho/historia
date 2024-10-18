package middleware

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/daussho/historia/domain/common"
	"github.com/daussho/historia/domain/user"
	"github.com/daussho/historia/utils/clock"
	context_util "github.com/daussho/historia/utils/context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"gorm.io/gorm"
)

type fiberHandler = func(ctx *fiber.Ctx) error

func AuthApi(db *gorm.DB) fiberHandler {
	return func(ctx *fiber.Ctx) error {
		log.Println("auth api")

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

		err := resolveUserSession(ctx, token, db)
		if err != nil {
			return common.NewErrorResponse(ctx, fiber.StatusUnauthorized, err)
		}

		return ctx.Next()
	}
}

func RateLimit() fiberHandler {
	cfg := limiter.ConfigDefault
	cfg.Max = 10
	cfg.Expiration = time.Second
	cfg.KeyGenerator = func(c *fiber.Ctx) string {
		return c.IP()
	}

	return limiter.New(cfg)
}

func AuthWeb(db *gorm.DB) fiberHandler {
	return func(ctx *fiber.Ctx) error {
		log.Println("auth web")

		token := ctx.Cookies("token")

		if token == "" {
			log.Println("token empty")
			return ctx.Redirect("/login")
		}

		err := resolveUserSession(ctx, token, db)
		if err != nil {
			log.Println("failed to resolve user session: ", err.Error())
			return ctx.Redirect("/login")
		}

		return ctx.Next()
	}
}

func resolveUserSession(ctx *fiber.Ctx, token string, db *gorm.DB) error {
	var userToken user.UserToken
	err := db.WithContext(ctx.Context()).First(&userToken, "token = ?", token).Error
	if err != nil {
		return fmt.Errorf("token not found")
	}

	if clock.Now().After(userToken.ExpiredAt) {
		return fmt.Errorf("token expired")
	}

	var userData user.User
	err = db.WithContext(ctx.Context()).First(&userData, "id = ?", userToken.UserID).Error
	if err != nil {
		return fmt.Errorf("user id %s not found", userToken.UserID)
	}

	sess := user.UserSession{
		User:      userData,
		ExpiredAt: userToken.ExpiredAt,
	}

	context_util.SetValue(ctx, common.UserSessionKey, sess)

	return nil
}
