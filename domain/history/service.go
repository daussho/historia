package history

import (
	"fmt"
	"time"

	"github.com/daussho/historia/domain/user"
	"github.com/daussho/historia/internal/trace"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	SaveVisit(c *fiber.Ctx, req VisitRequest) (string, error)
	UpdateVisit(ctx *fiber.Ctx, id string) error
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}

// Visit implements Service.
func (s *service) SaveVisit(ctx *fiber.Ctx, req VisitRequest) (string, error) {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "historyService.SaveVisit", nil)
	defer span.Finish()

	var userToken user.UserToken
	err := s.db.WithContext(ctx.Context()).First(&userToken, "token = ?", req.Token).Error
	if err != nil {
		return "", err
	}

	history := History{
		ID:           uuid.NewString(),
		Title:        req.Title,
		URL:          req.URL,
		UserID:       userToken.UserID,
		DeviceName:   req.DeviceName,
		LastActiveAt: time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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
	err := s.db.First(&history, "id = ?", id).Error
	if err != nil {
		return fmt.Errorf("history id %s not found", id)
	}

	history.LastActiveAt = time.Now()
	history.UpdatedAt = time.Now()

	err = s.db.Save(&history).Error
	if err != nil {
		return fmt.Errorf("failed to update history: %w", err)
	}

	return nil
}
