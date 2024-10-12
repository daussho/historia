package history

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service interface {
	SaveVisit(c *fiber.Ctx, req VisitRequest) error
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
func (s *service) SaveVisit(c *fiber.Ctx, req VisitRequest) error {
	history := History{
		Title:      req.Title,
		URL:        req.URL,
		UserID:     req.UserID,
		DeviceName: req.DeviceName,
		CreatedAt:  time.Now().UnixMilli(),
		UpdatedAt:  time.Now().UnixMilli(),
	}

	err := s.db.Create(&history).Error
	if err != nil {
		return err
	}

	return nil
}
