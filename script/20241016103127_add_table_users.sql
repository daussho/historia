-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id varchar(36) NOT NULL,
  name text DEFAULT NULL,
  email varchar(255) DEFAULT NULL UNIQUE,
  password varchar(255) DEFAULT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT NULL,
  PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
