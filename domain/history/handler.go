package history

import (
	"encoding/json"

	"github.com/daussho/historia/domain/common"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	SaveVisit(c *fiber.Ctx) error
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

	h.service.SaveVisit(ctx, req)

	return nil
}
