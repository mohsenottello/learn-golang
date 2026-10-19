-- +goose Up
ALTER TABLE resumes
    ADD COLUMN side_bg     TEXT NOT NULL DEFAULT '#7f7a72',
    ADD COLUMN side_text   TEXT NOT NULL DEFAULT '#ffffff',
    ADD COLUMN main_bg     TEXT NOT NULL DEFAULT '#ffffff',
    ADD COLUMN main_text   TEXT NOT NULL DEFAULT '#222222',
    ADD COLUMN accent      TEXT NOT NULL DEFAULT '#555555';   -- borders / section rules

-- +goose Down
ALTER TABLE resumes DROP COLUMN side_bg, DROP COLUMN side_text, DROP COLUMN main_bg, DROP COLUMN main_text, DROP COLUMN accent;
