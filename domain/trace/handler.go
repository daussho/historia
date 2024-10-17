package trace

import (
	"encoding/json"

	"github.com/daussho/historia/domain/common"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler interface {
}

type handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) Handler {
	return &handler{
		db,
	}
}

func (h *handler) SaveTrace(ctx *fiber.Ctx) error {
	var req []Trace
	err := json.Unmarshal(ctx.Body(), &req)
	if err != nil {
		return common.NewErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(&req)
	if err != nil {
		return common.NewErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	err = h.db.Create(&req).Error
	if err != nil {
		return common.NewErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return common.NewResponse(ctx, "success", nil)
}
