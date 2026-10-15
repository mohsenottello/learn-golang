-- +goose Up
CREATE TABLE admins (
    id            BIGSERIAL PRIMARY KEY,
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    verified      BOOLEAN NOT NULL DEFAULT FALSE,
    verify_token  TEXT UNIQUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    token      TEXT PRIMARY KEY,
    admin_id   BIGINT NOT NULL REFERENCES admins(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX sessions_expires_at_idx ON sessions (expires_at);

-- +goose Down
DROP TABLE sessions;
DROP TABLE admins;
