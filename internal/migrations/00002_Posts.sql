-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd  

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;

-- +goose StatementEnd