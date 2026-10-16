-- +goose Up
ALTER TABLE users
    ADD COLUMN password_hash    TEXT NOT NULL DEFAULT '',
    ADD COLUMN verified         BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN verify_token     TEXT UNIQUE,
    ADD COLUMN reset_token      TEXT UNIQUE,
    ADD COLUMN reset_expires_at TIMESTAMPTZ;

CREATE TABLE user_sessions (
    token      TEXT PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX user_sessions_user_id_idx ON user_sessions (user_id);

-- +goose Down
DROP TABLE user_sessions;
ALTER TABLE users
    DROP COLUMN password_hash,
    DROP COLUMN verified,
    DROP COLUMN verify_token,
    DROP COLUMN reset_token,
    DROP COLUMN reset_expires_at;
