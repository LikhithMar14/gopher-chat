-- +goose Up
-- +goose StatementBegin
ALTER TABLE posts ADD COLUMN tags TEXT[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE posts DROP COLUMN tags;
-- +goose StatementEnd
