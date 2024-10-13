package history

import (
	"time"

	"github.com/daussho/historia/domain/user"
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
		LastActiveAt: time.Now().UnixMilli(),
		CreatedAt:    time.Now().UnixMilli(),
		UpdatedAt:    time.Now().UnixMilli(),
	}

	res := s.db.WithContext(ctx.Context()).Create(&history)
	if res.Error != nil {
		return "", res.Error
	}

	return history.ID, nil
}

func (s *service) UpdateVisit(ctx *fiber.Ctx, id string) error {
	// var userToken user.UserToken
	// err := s.db.WithContext(ctx.Context()).First(&userToken, "token = ?", req.Token).Error
	// if err != nil {
	// 	return err
	// }

	err := s.db.Model(&History{}).
		Where("id = ?", id).
		Update("last_active_at", time.Now().UnixMilli()).
		Update("updated_at", time.Now().UnixMilli()).
		Error
	if err != nil {
		return err
	}

	return nil
}
