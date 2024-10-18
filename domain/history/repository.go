package history

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/daussho/historia/internal/trace"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/gofiber/fiber/v2"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

type GetPaginatedRequest struct {
	UserID string
	Offset uint64
	Limit  uint64
}

func (r *repository) GetPaginated(ctx *fiber.Ctx, req GetPaginatedRequest) ([]History, error) {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "historyRepository.SaveVisit", nil)
	defer span.Finish()

	var model History

	q := squirrel.Select(model.Columns()...).From(model.TableName())

	if req.UserID != "" {
		q = q.Where("user_id = ?", req.UserID)
	}

	q = q.Offset(req.Offset).Limit(req.Limit)

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var res []History
	err = sqlscan.Select(ctx.Context(), r.db, &res, sql, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
