package common

import "github.com/gofiber/fiber/v2"

type ContextKey string

const (
	UserContextKey ContextKey = "user"
)

type Response struct {
	Message string  `json:"message"`
	Error   *string `json:"error"`
	Data    any     `json:"data"`
}

func NewResponse(ctx *fiber.Ctx, message string, data any) error {
	return ctx.JSON(&Response{
		Message: message,
		Data:    data,
	})
}

func NewErrorResponse(ctx *fiber.Ctx, status int, err error) error {
	msg := err.Error()
	ctx.JSON(&Response{
		Error: &msg,
	})

	return ctx.SendStatus(status)
}
