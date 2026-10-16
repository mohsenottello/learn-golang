-- +goose Up
ALTER TABLE resumes
    ADD COLUMN default_lang TEXT NOT NULL DEFAULT 'de',
    ADD COLUMN title_en     TEXT NOT NULL DEFAULT '',
    ADD COLUMN subtitle_en  TEXT NOT NULL DEFAULT '',
    ADD COLUMN about_en     TEXT NOT NULL DEFAULT '';

ALTER TABLE resume_languages
    ADD COLUMN name_en  TEXT NOT NULL DEFAULT '',
    ADD COLUMN level_en TEXT NOT NULL DEFAULT '';

ALTER TABLE resume_education
    ADD COLUMN degree_en TEXT NOT NULL DEFAULT '',
    ADD COLUMN school_en TEXT NOT NULL DEFAULT '';

ALTER TABLE resume_experience
    ADD COLUMN role_en     TEXT NOT NULL DEFAULT '',
    ADD COLUMN subtitle_en TEXT NOT NULL DEFAULT '',
    ADD COLUMN location_en TEXT NOT NULL DEFAULT '';

ALTER TABLE resume_experience_bullets
    ADD COLUMN text_en TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE resume_experience_bullets DROP COLUMN text_en;
ALTER TABLE resume_experience DROP COLUMN role_en, DROP COLUMN subtitle_en, DROP COLUMN location_en;
ALTER TABLE resume_education DROP COLUMN degree_en, DROP COLUMN school_en;
ALTER TABLE resume_languages DROP COLUMN name_en, DROP COLUMN level_en;
ALTER TABLE resumes DROP COLUMN default_lang, DROP COLUMN title_en, DROP COLUMN subtitle_en, DROP COLUMN about_en;
