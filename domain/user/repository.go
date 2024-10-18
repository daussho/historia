package user

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/daussho/historia/internal/trace"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetUserByEmail(ctx *fiber.Ctx, email string) (User, error) {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "userRepository.GetUserByEmail", nil)
	defer span.Finish()

	var user User

	sql, args, err := sq.Select(user.Columns()...).From(user.TableName()).Where(sq.Eq{"email": email}).ToSql()
	if err != nil {
		return user, err
	}

	err = r.db.GetContext(ctx.Context(), &user, sql, args...)

	return user, err
}

func (r *repository) InsertToken(ctx *fiber.Ctx, req UserToken) error {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "userRepository.InsertToken", nil)
	defer span.Finish()

	now := sq.Expr("now()")
	expired := sq.Expr("now() + INTERVAL 7 day")
	data := map[string]any{
		"user_id":    req.UserID,
		"token":      req.Token,
		"expired_at": expired,
		"created_at": now,
		"updated_at": now,
	}

	sql, args, err := sq.Insert(req.TableName()).
		SetMap(data).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx.Context(), sql, args...)
	if err != nil {
		return err
	}

	return nil
}
