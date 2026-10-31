-- +goose Up
ALTER TABLE site_settings
    ADD COLUMN tagline_en      TEXT NOT NULL DEFAULT '',
    ADD COLUMN contact_text_en TEXT NOT NULL DEFAULT '',
    ADD COLUMN default_lang    TEXT NOT NULL DEFAULT 'de' CHECK (default_lang IN ('de', 'en'));

-- home becomes {"de": {...}, "en": {...}}; wrap the existing single-language document as "de".
UPDATE site_settings
SET home = jsonb_build_object('de', home)
WHERE home <> '{}'::jsonb AND NOT (home ? 'de') AND NOT (home ? 'en');

-- +goose Down
UPDATE site_settings SET home = COALESCE(home -> 'de', '{}'::jsonb) WHERE home ? 'de';
ALTER TABLE site_settings DROP COLUMN tagline_en, DROP COLUMN contact_text_en, DROP COLUMN default_lang;
