-- +goose Up
-- +goose StatementBegin
CREATE TABLE histories (
  id varchar(36) NOT NULL PRIMARY KEY,
  title longtext NOT NULL,
  url longtext NOT NULL,
  user_id longtext NOT NULL,
  device_name longtext NOT NULL,
  last_active_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS histories;
-- +goose StatementEnd
