-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_comments_content ON comments USING gin(content gin_trgm_ops);

CREATE INDEX idx_posts_title ON posts USING gin(title gin_trgm_ops);
CREATE INDEX idx_posts_tags ON posts USING gin(tags) ;

CREATE INDEX idx_users_username ON users USING gin(username gin_trgm_ops);

CREATE INDEX idx_posts_user_id ON posts (user_id);

CREATE INDEX idx_comments_post_id ON comments (post_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_comments_content;

DROP INDEX idx_posts_title;

DROP INDEX idx_posts_tags;

DROP INDEX idx_users_username;

DROP INDEX idx_posts_user_id;

DROP INDEX idx_comments_post_id;

-- +goose StatementEnd