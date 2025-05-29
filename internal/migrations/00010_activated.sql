-- +goose Up
-- +goose StatementBegin

    ALTER TABLE users ADD COLUMN IF NOT EXISTS activated BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

    ALTER TABLE users DROP COLUMN IF EXISTS activated;
-- +goose StatementEnd