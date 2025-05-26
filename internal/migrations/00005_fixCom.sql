-- +goose Up
-- +goose StatementBegin

-- Step 1: Drop default sequences attached to post_id and user_id (created by BIGSERIAL)
ALTER TABLE comments ALTER COLUMN post_id DROP DEFAULT;
ALTER TABLE comments ALTER COLUMN user_id DROP DEFAULT;

-- Step 2: Change column types from BIGSERIAL to BIGINT
-- Note: BIGSERIAL is shorthand for BIGINT + sequence, so just changing to BIGINT
ALTER TABLE comments ALTER COLUMN post_id TYPE BIGINT;
ALTER TABLE comments ALTER COLUMN user_id TYPE BIGINT;

-- Step 3: Add foreign key constraints
ALTER TABLE comments
    ADD CONSTRAINT fk_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;

ALTER TABLE comments
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Step 1: Drop foreign key constraints
ALTER TABLE comments DROP CONSTRAINT IF EXISTS fk_post;
ALTER TABLE comments DROP CONSTRAINT IF EXISTS fk_user;

-- Step 2: Revert columns to BIGSERIAL (creates new sequences)
ALTER TABLE comments ALTER COLUMN post_id TYPE BIGSERIAL;
ALTER TABLE comments ALTER COLUMN user_id TYPE BIGSERIAL;

-- +goose StatementEnd
