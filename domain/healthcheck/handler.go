package healthcheck

import (
	"fmt"

	"github.com/daussho/historia/internal/logger"
	"github.com/daussho/historia/internal/ntfy"
	"github.com/daussho/historia/internal/trace"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type Handler interface {
	Healthcheck(ctx *fiber.Ctx) error
}

type handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) Handler {
	return &handler{
		db,
	}
}

func (h *handler) Healthcheck(ctx *fiber.Ctx) error {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "healthcheckHandler.Healthcheck", nil)
	defer span.Finish()

	var res string
	err := h.db.GetContext(ctx.Context(), &res, "SELECT 1;")
	if err != nil {
		logger.Log().Errorf("healthcheck failed, err: %s", err.Error())
		go ntfy.SendError("Healthcheck error", fmt.Sprintf("healthcheck failed: %s", err.Error()))
		return ctx.SendStatus(500)
	}

	ctx.SendString("OK")
	return ctx.SendStatus(200)
}
