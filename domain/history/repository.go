package history

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/daussho/historia/internal/logger"
	"github.com/daussho/historia/internal/trace"
	"github.com/daussho/historia/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

type GetPaginatedRequest struct {
	UserID string
	Offset uint64
	Limit  uint64
}

func (r *repository) GetPaginated(ctx *fiber.Ctx, req GetPaginatedRequest) ([]History, error) {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "historyRepository.GetPaginated", nil)
	defer span.Finish()

	var model History

	q := sq.Select(model.Columns()...).
		From(model.TableName()).
		Offset(req.Offset).
		Limit(req.Limit).
		OrderBy("created_at DESC")

	if req.UserID != "" {
		q = q.Where("user_id = ?", req.UserID)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	logger.Debug().Msgf("%s; %s", sql, utils.JsonStringify(args))

	var res []History

	err = r.db.SelectContext(ctx.Context(), &res, sql, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *repository) SaveVisit(ctx *fiber.Ctx, req VisitRequest) (string, error) {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "historyRepository.SaveVisit", nil)
	defer span.Finish()

	id := uuid.NewString()
	now := sq.Expr("now()")

	data := map[string]any{
		"id":             id,
		"title":          req.Title,
		"url":            req.URL,
		"user_id":        req.UserID,
		"device_name":    req.DeviceName,
		"last_active_at": now,
		"created_at":     now,
		"updated_at":     now,
	}

	sql, args, err := sq.Insert(History{}.TableName()).SetMap(data).ToSql()
	if err != nil {
		return "", err
	}

	logger.Debug().Msgf("%s; %s", sql, utils.JsonStringify(args))
	_, err = r.db.ExecContext(ctx.Context(), sql, args...)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *repository) UpdateVisit(ctx *fiber.Ctx, id string) error {
	span, ctx := trace.StartSpanWithFiberCtx(ctx, "historyRepository.UpdateVisit", nil)
	defer span.Finish()

	now := sq.Expr("now()")

	data := map[string]any{
		"last_active_at": now,
		"updated_at":     now,
	}

	sql, args, err := sq.Update(History{}.TableName()).
		Where("id = ?", id).
		SetMap(data).
		ToSql()
	if err != nil {
		return err
	}

	logger.Debug().Msgf("%s; %s", sql, utils.JsonStringify(args))

	_, err = r.db.ExecContext(ctx.Context(), sql, args...)
	if err != nil {
		return err
	}

	return nil
}
