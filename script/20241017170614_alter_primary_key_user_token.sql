-- +goose Up
-- +goose StatementBegin
alter table user_tokens drop primary key, add primary key(user_id, token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table user_tokens drop primary key, add primary key(token);
-- +goose StatementEnd
