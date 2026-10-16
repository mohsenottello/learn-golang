-- +goose Up
-- Responsibilities become their own table.
CREATE TABLE resume_experience_bullets (
    id            BIGSERIAL PRIMARY KEY,
    experience_id BIGINT NOT NULL REFERENCES resume_experience(id) ON DELETE CASCADE,
    text          TEXT NOT NULL,
    position      INT NOT NULL DEFAULT 0
);
CREATE INDEX resume_experience_bullets_exp_idx ON resume_experience_bullets (experience_id, position);

INSERT INTO resume_experience_bullets (experience_id, text, position)
SELECT e.id, b.text, b.ord - 1
FROM resume_experience e, unnest(e.bullets) WITH ORDINALITY AS b(text, ord)
WHERE b.text <> '';

ALTER TABLE resume_experience DROP COLUMN bullets;

-- Dates become real DATE columns (month precision; NULL date_to = present).
-- Previously free text, so they cannot be converted and are reset.
ALTER TABLE resume_experience
    DROP COLUMN date_from, DROP COLUMN date_to,
    ADD COLUMN date_from DATE, ADD COLUMN date_to DATE;
ALTER TABLE resume_education
    DROP COLUMN date_from, DROP COLUMN date_to,
    ADD COLUMN date_from DATE, ADD COLUMN date_to DATE;

-- +goose Down
ALTER TABLE resume_education
    DROP COLUMN date_from, DROP COLUMN date_to,
    ADD COLUMN date_from TEXT NOT NULL DEFAULT '', ADD COLUMN date_to TEXT NOT NULL DEFAULT '';
ALTER TABLE resume_experience
    DROP COLUMN date_from, DROP COLUMN date_to,
    ADD COLUMN date_from TEXT NOT NULL DEFAULT '', ADD COLUMN date_to TEXT NOT NULL DEFAULT '',
    ADD COLUMN bullets TEXT[] NOT NULL DEFAULT '{}';
DROP TABLE resume_experience_bullets;
