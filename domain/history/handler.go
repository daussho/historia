package history

import (
	"encoding/json"
	"fmt"

	"github.com/daussho/historia/domain/common"
	"github.com/daussho/historia/internal/trace"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	SaveVisit(c *fiber.Ctx) error
	UpdateVisit(ctx *fiber.Ctx) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) SaveVisit(ctx *fiber.Ctx) error {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "historyHandler.SaveVisit", nil)
	defer span.Finish()

	var req VisitRequest
	err := json.Unmarshal(ctx.Body(), &req)
	if err != nil {
		return common.NewErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(&req)
	if err != nil {
		return common.NewErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	id, err := h.service.SaveVisit(ctx, req)
	if err != nil {
		return common.NewErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return common.NewResponse(ctx, "success", map[string]any{
		"id": id,
	})
}

func (h *handler) UpdateVisit(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return common.NewErrorResponse(ctx, fiber.StatusBadRequest, fmt.Errorf("invalid id"))
	}

	err := h.service.UpdateVisit(ctx, id)
	if err != nil {
		return common.NewErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return common.NewResponse(ctx, "success", nil)
}
