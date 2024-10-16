package healthcheck

import (
	"github.com/daussho/historia/internal/trace"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler interface {
	Healthcheck(ctx *fiber.Ctx) error
}

type handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) Handler {
	return &handler{
		db,
	}
}

func (h *handler) Healthcheck(ctx *fiber.Ctx) error {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "healthcheckHandler.Healthcheck", nil)
	defer span.Finish()

	err := h.db.Raw("SELECT 1").Error
	if err != nil {
		return ctx.SendStatus(500)
	}

	ctx.SendString("OK")
	return ctx.SendStatus(200)
}