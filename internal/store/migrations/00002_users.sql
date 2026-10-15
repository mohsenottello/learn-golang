-- +goose Up
CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    email      TEXT UNIQUE NOT NULL,
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE users;
