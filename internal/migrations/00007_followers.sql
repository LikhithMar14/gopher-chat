-- +goose Up
CREATE TABLE IF NOT EXISTS followers (
    user_id BIGINT NOT NULL,
    follower_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, follower_id)
);

-- +goose Down
DROP TABLE IF EXISTS followers;