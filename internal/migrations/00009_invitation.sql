-- +goose Up
-- +goose StatementBegin

        CREATE TABLE IF NOT EXISTS user_invitations (
            token bytea NOT NULL,
            user_id bigint NOT NULL,
            PRIMARY KEY (token , user_id)
        );
        CREATE INDEX IF NOT EXISTS idx_user_invitations_user_id ON user_invitations (user_id);
        CREATE INDEX IF NOT EXISTS idx_user_invitations_token ON user_invitations (token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

        DROP TABLE IF EXISTS user_invitations;
-- +goose StatementEnd
