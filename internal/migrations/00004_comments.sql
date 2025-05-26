-- +goose Up
-- +goose StatementBegin
CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGSERIAL NOT NULL,
    user_id BIGSERIAL NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE comments;
-- +goose StatementEnd