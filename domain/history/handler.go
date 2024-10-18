package history

import (
	"encoding/json"
	"fmt"

	"github.com/daussho/historia/domain/common"
	"github.com/daussho/historia/domain/user"
	"github.com/daussho/historia/internal/trace"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	SaveVisit(c *fiber.Ctx) error
	UpdateVisit(ctx *fiber.Ctx) error
	ListHistory(ctx *fiber.Ctx) error
}

type handler struct {
	service Service
	userSvc user.Service
}

func NewHandler(service Service, userSvc user.Service) Handler {
	return &handler{
		service: service,
		userSvc: userSvc,
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

	user, ok := h.userSvc.GetSession(ctx)
	if !ok {
		return common.NewErrorResponse(ctx, fiber.StatusUnauthorized, fmt.Errorf("unauthorized"))
	}

	id, err := h.service.SaveVisit(ctx, req, user.ID)
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

func (h *handler) ListHistory(ctx *fiber.Ctx) error {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "historyHandler.ListHistory", nil)
	defer span.Finish()

	user, ok := h.userSvc.GetSession(ctx)
	if !ok {
		return common.NewErrorResponse(ctx, fiber.StatusUnauthorized, fmt.Errorf("unauthorized"))
	}

	pageSize := ctx.QueryInt("page_size", 50)
	pageIndex := ctx.QueryInt("page_index", 1)

	histories, err := h.service.ListHistory(ctx, user.ID, pageSize, pageIndex)
	if err != nil {
		return ctx.Redirect("/history")
	}

	nextUrl := fmt.Sprintf("/history?page_index=%d", pageIndex+1)
	prevUrl := fmt.Sprintf("/history?page_index=%d", pageIndex-1)
	if pageIndex == 1 {
		prevUrl = "#"
	}

	type historyData map[string]string
	data := make([]historyData, 0, len(histories))
	for _, history := range histories {
		data = append(data, historyData{
			"title":      history.Title,
			"url":        history.URL,
			"created_at": history.CreatedAt.Format("2006-01-02 15:04:05"),
			"duration":   history.LastActiveAt.Sub(history.CreatedAt).String(),
		})
	}

	return ctx.Render("page/history", fiber.Map{
		"histories": data,
		"next_url":  nextUrl,
		"prev_url":  prevUrl,
	})
}
