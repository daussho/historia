package user

import (
	"fmt"
	"log"

	"github.com/daussho/historia/domain/common"
	"github.com/daussho/historia/internal/trace"
	context_util "github.com/daussho/historia/utils/context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Login(ctx *fiber.Ctx) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Login(ctx *fiber.Ctx) error {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "userHandler.Login", nil)
	defer span.Finish()

	_, ok := context_util.GetValue(ctx, common.UserSessionKey).(UserSession)
	if ok {
		return ctx.Redirect("/history")
	}

	if ctx.Method() == "GET" {
		return h.renderLogin(ctx)
	} else if ctx.Method() == "POST" {
		return h.validateLogin(ctx)
	} else {
		return common.NewErrorResponse(ctx, fiber.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
	}
}

func (h *handler) renderLogin(ctx *fiber.Ctx) error {
	return ctx.Render("page/login", fiber.Map{
		"err_msg": ctx.Query("err_msg"),
	})
}

func (h *handler) validateLogin(ctx *fiber.Ctx) error {
	var req LoginRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("failed to parse body: ", err.Error())
		return ctx.Redirect("/login?err_msg=invalid+request")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(&req)
	if err != nil {
		log.Println("failed to validate request: ", err.Error())
		return ctx.Redirect("/login?err_msg=invalid+request")
	}

	token, err := h.service.GenerateToken(ctx, req)
	if err != nil {
		return ctx.Redirect("/login?err_msg=invalid+username+or+password")
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token.Token,
		Expires:  token.ExpiredAt,
		HTTPOnly: true,
	})

	return ctx.Redirect("/history")
}
