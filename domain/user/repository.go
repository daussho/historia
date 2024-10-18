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
