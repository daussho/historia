-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS histories_user_id ON histories (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS histories_user_id;

-- +goose StatementEnd
