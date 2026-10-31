package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

// Langs the site content can be stored in.
var Langs = []string{"de", "en"}

type Site struct {
	Name          string                 `json:"name"`
	Tagline       string                 `json:"tagline"`
	TaglineEn     string                 `json:"tagline_en"`
	ContactEmail  string                 `json:"contact_email"`
	ContactText   string                 `json:"contact_text"`
	ContactTextEn string                 `json:"contact_text_en"`
	PrimaryColor  string                 `json:"primary_color"`
	Registration  bool                   `json:"registration"`
	DefaultLang   string                 `json:"default_lang"`
	HasLogo       bool                   `json:"has_logo"`
	HasPortrait   bool                   `json:"has_portrait"`
	Home          map[string]HomeContent `json:"home"` // keyed by lang: "de", "en"
	UpdatedAt     time.Time              `json:"updated_at"`
}

type Feature struct {
	Icon  string `json:"icon"`  // material icon name
	Color string `json:"color"` // teal | purple | pink | blue
	Title string `json:"title"`
	Text  string `json:"text"`
}

type HomeContent struct {
	Greeting        string    `json:"greeting"`
	Headline        string    `json:"headline"`
	Description     string    `json:"description"`
	PrimaryLabel    string    `json:"primary_label"`
	SecondaryLabel  string    `json:"secondary_label"`
	Note            string    `json:"note"`
	FloatingNote    string    `json:"floating_note"`
	FeaturedSlug    string    `json:"featured_slug"` // resume that "My Resume" opens
	FeaturesLabel   string    `json:"features_label"`
	FeaturesHeading string    `json:"features_heading"`
	FeaturesIntro   string    `json:"features_intro"`
	Features        []Feature `json:"features"`
	AboutLabel      string    `json:"about_label"`
	AboutHeading    string    `json:"about_heading"`
	AboutText       string    `json:"about_text"`
	AboutButton     string    `json:"about_button"`
	CTAHeading      string    `json:"cta_heading"`
	CTAText         string    `json:"cta_text"`
	CTAButton       string    `json:"cta_button"`
	FooterText      string    `json:"footer_text"`
}

// DefaultHome is used for a language that has nothing saved yet (and as the "reset" target).
func DefaultHome(lang string) HomeContent {
	if lang == "de" {
		return HomeContent{
			Greeting: "Hallo, ich bin", Headline: "Elahe Roshandel",
			Description:    "Willkommen auf meiner persönlichen Seite! Hier teile ich meinen beruflichen Werdegang, meine Erfahrung und meine Fähigkeiten. Du kannst hier auch deinen eigenen Lebenslauf erstellen und teilen.",
			PrimaryLabel:   "Mein Lebenslauf", SecondaryLabel: "Eigenen Lebenslauf erstellen",
			Note:           "Deine Geschichte verdient eine professionelle Präsentation.",
			FloatingNote:   "Deine Zukunft, deine Geschichte.",
			FeaturedSlug:   "elahe-roshandel",
			FeaturesLabel:  "Was du hier tun kannst", FeaturesHeading: "Lebenslauf erstellen und teilen",
			FeaturesIntro:  "Ob du einen neuen Job suchst, bessere Chancen erkunden oder deine Fähigkeiten zeigen möchtest – hier ist ein guter Startpunkt.",
			Features: []Feature{
				{"description", "teal", "Lebenslauf erstellen", "Erstelle einen professionellen Lebenslauf mit einem einfachen, geführten Prozess."},
				{"edit", "purple", "Mach ihn zu deinem", "Passe Erfahrung, Ausbildung, Fähigkeiten und Darstellung an."},
				{"share", "pink", "Selbstbewusst teilen", "Teile deinen Lebenslauf mit Recruitern, Arbeitgebern und Kollegen."},
				{"star_outline", "blue", "Hervorstechen", "Präsentiere deinen beruflichen Hintergrund in einem klaren, polierten Format."},
			},
			AboutLabel: "Mein Lebenslauf", AboutHeading: "Lern mich kennen",
			AboutText:   "Neugierig auf meinen beruflichen Weg? Auf meinem Lebenslauf erfährst du mehr über Hintergrund, Erfahrung, Ausbildung und Fähigkeiten.",
			AboutButton: "Lebenslauf ansehen",
			CTAHeading:  "Deine nächste Chance beginnt hier",
			CTAText:     "Erstelle einen Lebenslauf, der deine Fähigkeiten, Erfahrungen und Ziele widerspiegelt. Mach ihn zu deinem und erzähle deine Geschichte mit Selbstvertrauen.",
			CTAButton:   "Lebenslauf erstellen",
			FooterText:  "Mein persönliches Portfolio und ein Ort, um deinen eigenen Lebenslauf zu erstellen und zu teilen.",
		}
	}
	return HomeContent{
		Greeting: "Hello, I'm", Headline: "Elahe Roshandel",
		Description:    "Welcome to my personal space! Here, I share my professional journey, experience, and skills. You can also create your own professional resume and share it with the world.",
		PrimaryLabel:   "View My Resume", SecondaryLabel: "Create Your Resume",
		Note:           "Your story deserves a professional presentation.", FloatingNote: "Your future, your story.",
		FeaturedSlug:   "elahe-roshandel",
		FeaturesLabel:  "What You Can Do", FeaturesHeading: "Create and Share Your Resume",
		FeaturesIntro:  "Whether you're looking for a new job, exploring better opportunities, or showcasing your skills, this is a great place to start.",
		Features: []Feature{
			{"description", "teal", "Create Your Resume", "Build a professional resume with a simple, guided process."},
			{"edit", "purple", "Make It Yours", "Personalize your experience, education, skills, and presentation."},
			{"share", "pink", "Share with Confidence", "Share your resume with recruiters, employers, and colleagues."},
			{"star_outline", "blue", "Stand Out", "Present your professional background in a clean, polished format."},
		},
		AboutLabel: "My Resume", AboutHeading: "Get to Know Me",
		AboutText:   "Curious about my professional journey? Visit my resume to learn more about my background, experience, education, and skills.",
		AboutButton: "Explore My Resume",
		CTAHeading:  "Your Next Opportunity Starts Here",
		CTAText:     "Create a resume that reflects your skills, experience, and ambitions. Make it yours and share your story with confidence.",
		CTAButton:   "Create Your Resume",
		FooterText:  "My personal portfolio and a place to create and share your professional resume.",
	}
}

// parseHome reads the JSONB column into a per-language map. It also accepts the
// pre-migration single-document shape (treated as "de") so a half-migrated row still loads.
func parseHome(raw []byte) (map[string]HomeContent, error) {
	out := map[string]HomeContent{}
	if len(raw) == 0 {
		return out, nil
	}
	var byLang map[string]json.RawMessage
	if err := json.Unmarshal(raw, &byLang); err != nil {
		return nil, err
	}
	if _, legacy := byLang["headline"]; legacy {
		var h HomeContent
		if err := json.Unmarshal(raw, &h); err != nil {
			return nil, err
		}
		out["de"] = h
		return out, nil
	}
	for lang, v := range byLang {
		var h HomeContent
		if err := json.Unmarshal(v, &h); err != nil {
			return nil, err
		}
		out[lang] = h
	}
	return out, nil
}

func (s *Store) GetSite(ctx context.Context) (*Site, error) {
	var x Site
	var home []byte
	err := s.db.QueryRowContext(ctx, `
		SELECT name, tagline, tagline_en, contact_email, contact_text, contact_text_en, primary_color,
		       logo IS NOT NULL, portrait IS NOT NULL, home, registration, default_lang, updated_at
		FROM site_settings WHERE id = 1`,
	).Scan(&x.Name, &x.Tagline, &x.TaglineEn, &x.ContactEmail, &x.ContactText, &x.ContactTextEn, &x.PrimaryColor,
		&x.HasLogo, &x.HasPortrait, &home, &x.Registration, &x.DefaultLang, &x.UpdatedAt)
	if err != nil {
		return nil, translate(err)
	}
	if x.Home, err = parseHome(home); err != nil {
		return nil, err
	}
	for _, lang := range Langs { // a language with nothing saved gets the defaults
		if x.Home[lang].Headline == "" {
			x.Home[lang] = DefaultHome(lang)
		}
	}
	return &x, nil
}

func (s *Store) UpdateSite(ctx context.Context, x *Site) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE site_settings
		SET name=$1, tagline=$2, tagline_en=$3, contact_email=$4, contact_text=$5, contact_text_en=$6,
		    primary_color=$7, registration=$8, default_lang=$9, updated_at=now()
		WHERE id = 1`,
		x.Name, x.Tagline, x.TaglineEn, x.ContactEmail, x.ContactText, x.ContactTextEn,
		x.PrimaryColor, x.Registration, x.DefaultLang)
	return err
}

// UpdateHome stores the home page content for one language.
// Pass the zero value to reset that language to its defaults.
func (s *Store) UpdateHome(ctx context.Context, lang string, h HomeContent) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var raw []byte
	if err := tx.QueryRowContext(ctx, `SELECT home FROM site_settings WHERE id = 1 FOR UPDATE`).Scan(&raw); err != nil {
		return translate(err)
	}
	all, err := parseHome(raw)
	if err != nil {
		return err
	}
	if h.Headline == "" {
		delete(all, lang) // reset → GetSite substitutes defaults
	} else {
		all[lang] = h
	}
	b, err := json.Marshal(all)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE site_settings SET home=$1::jsonb, updated_at=now() WHERE id = 1`, b); err != nil {
		return err
	}
	return tx.Commit()
}

// --- images: logo and portrait share one implementation ---

func (s *Store) setSiteImage(ctx context.Context, col string, data []byte, ct string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE site_settings SET `+col+`=$1, `+col+`_type=$2, updated_at=now() WHERE id = 1`, data, ct)
	return err
}

func (s *Store) siteImage(ctx context.Context, col string) ([]byte, string, error) {
	var data []byte
	var ct sql.NullString
	err := s.db.QueryRowContext(ctx,
		`SELECT `+col+`, `+col+`_type FROM site_settings WHERE id = 1 AND `+col+` IS NOT NULL`).Scan(&data, &ct)
	return data, ct.String, translate(err)
}

func (s *Store) SetSiteLogo(ctx context.Context, data []byte, ct string) error {
	return s.setSiteImage(ctx, "logo", data, ct)
}
func (s *Store) SiteLogo(ctx context.Context) ([]byte, string, error) {
	return s.siteImage(ctx, "logo")
}
func (s *Store) SetPortrait(ctx context.Context, data []byte, ct string) error {
	return s.setSiteImage(ctx, "portrait", data, ct)
}
func (s *Store) Portrait(ctx context.Context) ([]byte, string, error) {
	return s.siteImage(ctx, "portrait")
}
