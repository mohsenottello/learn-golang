-- +goose Up
ALTER TABLE site_settings
    ADD COLUMN home          JSONB NOT NULL DEFAULT '{}',
    ADD COLUMN primary_color TEXT  NOT NULL DEFAULT '#4f46e5',
    ADD COLUMN portrait      BYTEA,
    ADD COLUMN portrait_type TEXT;
-- +goose Down
ALTER TABLE site_settings DROP COLUMN home, DROP COLUMN primary_color, DROP COLUMN portrait, DROP COLUMN portrait_type;
