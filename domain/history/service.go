package history

import (
	"fmt"

	"github.com/daussho/historia/internal/trace"
	"github.com/gofiber/fiber/v2"
)

type Service interface {
	SaveVisit(c *fiber.Ctx, req VisitRequest, userID string) (string, error)
	UpdateVisit(ctx *fiber.Ctx, id string) error
	ListHistory(ctx *fiber.Ctx, userID string, pageSize, pageIndex int) ([]History, error)
}

type service struct {
	repo *repository
}

func NewService(repo *repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) SaveVisit(ctx *fiber.Ctx, req VisitRequest, userID string) (string, error) {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "historyService.SaveVisit", nil)
	defer span.Finish()

	req.UserID = userID

	id, err := s.repo.SaveVisit(ctx, req)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *service) UpdateVisit(ctx *fiber.Ctx, id string) error {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "historyService.UpdateVisit", nil)
	defer span.Finish()

	err := s.repo.UpdateVisit(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to update history: %w", err)
	}

	return nil
}

func (s *service) ListHistory(ctx *fiber.Ctx, userID string, pageSize, pageIndex int) ([]History, error) {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "historyService.ListHistory", nil)
	defer span.Finish()

	offset := 0
	if pageIndex > 0 {
		offset = pageSize * (pageIndex - 1)
	}

	req := GetPaginatedRequest{
		UserID: userID,
		Offset: uint64(offset),
		Limit:  uint64(pageSize),
	}
	histories, err := s.repo.GetPaginated(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list history: %w", err)
	}

	return histories, nil
}
