package history

import (
	"fmt"

	"github.com/daussho/historia/domain/user"
	"github.com/daussho/historia/internal/trace"
	"github.com/daussho/historia/utils/clock"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	SaveVisit(c *fiber.Ctx, req VisitRequest, userID string) (string, error)
	UpdateVisit(ctx *fiber.Ctx, id string) error
	ListHistory(ctx *fiber.Ctx, userID string, pageSize, pageIndex int) ([]History, error)
}

type service struct {
	db      *gorm.DB
	userSvc user.Service
}

func NewService(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}

// Visit implements Service.
func (s *service) SaveVisit(ctx *fiber.Ctx, req VisitRequest, userID string) (string, error) {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "historyService.SaveVisit", nil)
	defer span.Finish()

	history := History{
		ID:           uuid.NewString(),
		Title:        req.Title,
		URL:          req.URL,
		UserID:       userID,
		DeviceName:   req.DeviceName,
		LastActiveAt: clock.Now(),
		CreatedAt:    clock.Now(),
		UpdatedAt:    clock.Now(),
	}

	res := s.db.WithContext(ctx.Context()).Create(&history)
	if res.Error != nil {
		return "", res.Error
	}

	return history.ID, nil
}

func (s *service) UpdateVisit(ctx *fiber.Ctx, id string) error {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "historyService.UpdateVisit", nil)
	defer span.Finish()

	var history History
	err := s.db.WithContext(ctx.Context()).First(&history, "id = ?", id).Error
	if err != nil {
		return fmt.Errorf("history id %s not found", id)
	}

	history.LastActiveAt = clock.Now()
	history.UpdatedAt = clock.Now()

	err = s.db.WithContext(ctx.Context()).Save(&history).Error
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

	var histories []History
	err := s.db.WithContext(ctx.Context()).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&histories).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to list history: %w", err)
	}

	return histories, nil
}
