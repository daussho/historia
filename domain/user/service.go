package user

import (
	"fmt"
	"log"
	"time"

	"github.com/daussho/historia/domain/common"
	"github.com/daussho/historia/internal/trace"
	"github.com/daussho/historia/utils/clock"
	context_util "github.com/daussho/historia/utils/context"
	"github.com/daussho/historia/utils/password"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GenerateToken(ctx *fiber.Ctx, req LoginRequest) (UserToken, error)
	Login(ctx *fiber.Ctx, req LoginRequest) (User, error)
	GetSession(ctx *fiber.Ctx) (UserSession, bool)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}

func (s *service) GenerateToken(ctx *fiber.Ctx, req LoginRequest) (UserToken, error) {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "userService.GenerateToken", nil)
	defer span.Finish()

	var user User
	err := s.db.WithContext(ctx.Context()).First(&user, "email = ?", req.Email).Error
	if err != nil {
		log.Println("failed to get user: ", err.Error())
		return UserToken{}, fmt.Errorf("wrong email or password")
	}

	valid := password.Check(req.Password, user.Password)
	if !valid {
		log.Println("invalid password")
		return UserToken{}, fmt.Errorf("wrong email or password")
	}

	newToken := uuid.NewString()
	userToken := UserToken{
		UserID:    user.ID,
		Token:     newToken,
		ExpiredAt: clock.Now().Add(24 * time.Hour),
		CreatedAt: clock.Now(),
		UpdatedAt: clock.Now(),
	}
	err = s.db.WithContext(ctx.Context()).Create(&userToken).Error
	if err != nil {
		log.Println("failed to create user token: ", err.Error())
		return UserToken{}, fmt.Errorf("failed to create user token")
	}

	return userToken, nil
}

func (s *service) Login(ctx *fiber.Ctx, req LoginRequest) (User, error) {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "userService.Login", nil)
	defer span.Finish()

	var user User
	err := s.db.WithContext(ctx.Context()).First(&user, "email = ?", req.Email).Error
	if err != nil {
		log.Println("failed to get user: ", err.Error())
		return user, fmt.Errorf("wrong email or password")
	}

	valid := password.Check(req.Password, user.Password)
	if !valid {
		log.Println("invalid password")
		return user, fmt.Errorf("wrong email or password")
	}

	return user, nil
}

func (s *service) GetSession(ctx *fiber.Ctx) (UserSession, bool) {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "userService.GetSession", nil)
	defer span.Finish()

	session, ok := context_util.GetValue(ctx, common.UserSessionKey).(UserSession)
	if !ok {
		return UserSession{}, false
	}

	// check expired
	if clock.Now().After(session.ExpiredAt) {
		return UserSession{}, false
	}

	return session, true
}
