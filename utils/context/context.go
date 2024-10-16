package context_util

import (
	"github.com/daussho/historia/domain/user"
	"github.com/gofiber/fiber/v2"
)

func SetValue(ctx *fiber.Ctx, key string, value any) {
	ctx.Locals(key, value)
}

func GetValue(ctx *fiber.Ctx, key string) any {
	return ctx.Locals(key)
}

func SetUserCtx(ctx *fiber.Ctx, user user.User) {
	SetValue(ctx, "user", user)
}

func GetUserCtx(ctx *fiber.Ctx) user.User {
	userData, found := GetValue(ctx, "user").(user.User)
	if !found {
		return user.User{}
	}

	return userData
}
