package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/daussho/historia/domain/common"
	"github.com/daussho/historia/domain/user"
	context_util "github.com/daussho/historia/utils/context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"gorm.io/gorm"
)

type fiberHandler = func(ctx *fiber.Ctx) error

func AuthApi(db *gorm.DB) fiberHandler {
	return func(ctx *fiber.Ctx) error {
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

func RateLimit() fiberHandler {
	cfg := limiter.ConfigDefault
	cfg.Max = 10
	cfg.Expiration = time.Second
	cfg.KeyGenerator = func(c *fiber.Ctx) string {
		user := context_util.GetUserCtx(c)

		if user.ID == "" {
			return fmt.Sprintf("ip:%s", c.IP())
		}

		return fmt.Sprintf("user:%s", user.ID)
	}

	return limiter.New(cfg)
}
