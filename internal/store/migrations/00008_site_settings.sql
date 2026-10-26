-- +goose Up
CREATE TABLE site_settings (
    id            INT PRIMARY KEY CHECK (id = 1),   -- exactly one row, ever
    name          TEXT NOT NULL DEFAULT 'MyApp',
    tagline       TEXT NOT NULL DEFAULT '',
    contact_email TEXT NOT NULL DEFAULT '',
    contact_text  TEXT NOT NULL DEFAULT '',
    logo          BYTEA,
    logo_type     TEXT,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
INSERT INTO site_settings (id) VALUES (1);

-- +goose Down
DROP TABLE site_settings;
