-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_tokens (
  user_id varchar(36) NOT NULL PRIMARY KEY,
  token varchar(36) NOT NULL,
  expired_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_tokens;
-- +goose StatementEnd
