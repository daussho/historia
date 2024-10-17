-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS histories_created_at ON histories (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS histories_created_at;
-- +goose StatementEnd
