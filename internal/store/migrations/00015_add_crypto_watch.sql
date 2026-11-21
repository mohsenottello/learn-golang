-- +goose Up
ALTER TABLE users ADD COLUMN telegram_chat_id BIGINT, ADD COLUMN telegram_link_code TEXT;
CREATE TABLE crypto_watches (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    symbol     TEXT NOT NULL,
    pct        NUMERIC(6,2) NOT NULL,
    ref_price  NUMERIC(18,8) NOT NULL,
    ref_time   TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, symbol)
);
-- +goose Down
DROP TABLE crypto_watches;
ALTER TABLE users DROP COLUMN telegram_chat_id, DROP COLUMN telegram_link_code;
