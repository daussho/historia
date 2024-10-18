package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/daussho/historia/domain/common"
	"github.com/daussho/historia/domain/user"
	"github.com/daussho/historia/internal/logger"
	context_util "github.com/daussho/historia/utils/context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type fiberHandler = func(ctx *fiber.Ctx) error

type middleware struct {
	userSvc user.Service
}

func NewMiddleware(userSvc user.Service) *middleware {
	return &middleware{
		userSvc: userSvc,
	}
}

func (mw *middleware) AuthApi() fiberHandler {
	return func(ctx *fiber.Ctx) error {
		logger.Log().Debug("auth api")

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

		err := mw.resolveUserSession(ctx, token)
		if err != nil {
			return common.NewErrorResponse(ctx, fiber.StatusUnauthorized, err)
		}

		return ctx.Next()
	}
}

func (mw *middleware) RateLimit() fiberHandler {
	cfg := limiter.ConfigDefault
	cfg.Max = 10
	cfg.Expiration = time.Second
	cfg.KeyGenerator = func(c *fiber.Ctx) string {
		return c.IP()
	}

	return limiter.New(cfg)
}

func (mw *middleware) AuthWeb() fiberHandler {
	return func(ctx *fiber.Ctx) error {
		logger.Log().Debug("auth web")

		token := ctx.Cookies("token")

		if token == "" {
			logger.Log().Warn("token empty")
			return ctx.Redirect("/login")
		}

		err := mw.resolveUserSession(ctx, token)
		if err != nil {
			logger.Log().Warn("failed to resolve user session: %s", err.Error())
			return ctx.Redirect("/login")
		}

		return ctx.Next()
	}
}

func (mw *middleware) resolveUserSession(ctx *fiber.Ctx, token string) error {
	userSession, err := mw.userSvc.GetUserSessionWithToken(ctx, token)
	if err != nil {
		logger.Log().Warn("failed to get user session: %s", err.Error())
		return err
	}

	context_util.SetValue(ctx, common.UserSessionKey, userSession)

	return nil
}

func (mw *middleware) PanicRecovery() fiberHandler {
	return func(ctx *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				logger.Log().Error(fmt.Sprintf("%v", r))
				ctx.Status(fiber.StatusInternalServerError)
			}
		}()

		return ctx.Next()
	}
}
