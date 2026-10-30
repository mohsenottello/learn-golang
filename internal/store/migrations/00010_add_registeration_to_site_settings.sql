-- +goose Up
ALTER TABLE site_settings
    ADD COLUMN registration BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE site_settings DROP COLUMN registration;
