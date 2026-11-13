-- +goose Up
ALTER TABLE resumes ADD COLUMN access_hash TEXT NOT NULL DEFAULT '';
-- +goose Down
ALTER TABLE resumes DROP COLUMN access_hash;
