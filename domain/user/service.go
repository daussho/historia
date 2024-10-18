package user

import (
	"fmt"
	"time"

	"github.com/daussho/historia/domain/common"
	"github.com/daussho/historia/internal/logger"
	"github.com/daussho/historia/internal/trace"
	"github.com/daussho/historia/utils/clock"
	context_util "github.com/daussho/historia/utils/context"
	"github.com/daussho/historia/utils/password"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Service interface {
	GenerateToken(ctx *fiber.Ctx, req LoginRequest) (UserToken, error)
	Login(ctx *fiber.Ctx, req LoginRequest) (User, error)
	GetSession(ctx *fiber.Ctx) (UserSession, bool)
}

type service struct {
	repo *repository
}

func NewService(repo *repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GenerateToken(ctx *fiber.Ctx, req LoginRequest) (UserToken, error) {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "userService.GenerateToken", nil)
	defer span.Finish()

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.Log().Error(err)
		return UserToken{}, fmt.Errorf("wrong email or password")
	}

	valid := password.Check(req.Password, user.Password)
	if !valid {
		logger.Log().Error(err)
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

	err = s.repo.InsertToken(ctx, userToken)
	if err != nil {
		logger.Log().Error(err)
		return UserToken{}, fmt.Errorf("failed to create user token")
	}

	return userToken, nil
}

func (s *service) Login(ctx *fiber.Ctx, req LoginRequest) (User, error) {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "userService.Login", nil)
	defer span.Finish()

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.Log().Error(err)
		return user, fmt.Errorf("wrong email or password")
	}

	valid := password.Check(req.Password, user.Password)
	if !valid {
		logger.Log().Error("invalid password")
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
