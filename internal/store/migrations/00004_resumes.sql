-- +goose Up
CREATE TABLE resumes (
    user_id    BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    slug       TEXT UNIQUE NOT NULL,
    published  BOOLEAN NOT NULL DEFAULT FALSE,
    title      TEXT NOT NULL DEFAULT '',
    subtitle   TEXT NOT NULL DEFAULT '',
    about      TEXT NOT NULL DEFAULT '',
    email      TEXT NOT NULL DEFAULT '',
    phone      TEXT NOT NULL DEFAULT '',
    location   TEXT NOT NULL DEFAULT '',
    linkedin   TEXT NOT NULL DEFAULT '',
    photo      BYTEA,
    photo_type TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE resume_software (
    id       BIGSERIAL PRIMARY KEY,
    user_id  BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name     TEXT NOT NULL,
    position INT NOT NULL DEFAULT 0
);
CREATE INDEX resume_software_user_idx ON resume_software (user_id, position);

CREATE TABLE resume_languages (
    id       BIGSERIAL PRIMARY KEY,
    user_id  BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name     TEXT NOT NULL,
    level    TEXT NOT NULL DEFAULT '',
    position INT NOT NULL DEFAULT 0
);
CREATE INDEX resume_languages_user_idx ON resume_languages (user_id, position);

CREATE TABLE resume_education (
    id        BIGSERIAL PRIMARY KEY,
    user_id   BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    degree    TEXT NOT NULL,
    school    TEXT NOT NULL DEFAULT '',
    date_from TEXT NOT NULL DEFAULT '',
    date_to   TEXT NOT NULL DEFAULT '',
    position  INT NOT NULL DEFAULT 0
);
CREATE INDEX resume_education_user_idx ON resume_education (user_id, position);

CREATE TABLE resume_experience (
    id        BIGSERIAL PRIMARY KEY,
    user_id   BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company   TEXT NOT NULL,
    link      TEXT NOT NULL DEFAULT '',
    role      TEXT NOT NULL DEFAULT '',
    subtitle  TEXT NOT NULL DEFAULT '',
    date_from TEXT NOT NULL DEFAULT '',
    date_to   TEXT NOT NULL DEFAULT '',
    location  TEXT NOT NULL DEFAULT '',
    bullets   TEXT[] NOT NULL DEFAULT '{}',
    position  INT NOT NULL DEFAULT 0
);
CREATE INDEX resume_experience_user_idx ON resume_experience (user_id, position);

-- +goose Down
DROP TABLE resume_experience;
DROP TABLE resume_education;
DROP TABLE resume_languages;
DROP TABLE resume_software;
DROP TABLE resumes;
