// Package store is the only package that talks to the database.
// Each model gets its own file (admins.go, sessions.go, users.go, ...).
package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type Site struct {
	Name         string      `json:"name"`
	Tagline      string      `json:"tagline"`
	ContactEmail string      `json:"contact_email"`
	ContactText  string      `json:"contact_text"`
	PrimaryColor string      `json:"primary_color"`
	Registration bool        `json:"registration"`
	HasLogo      bool        `json:"has_logo"`
	HasPortrait  bool        `json:"has_portrait"`
	Home         HomeContent `json:"home"`
	UpdatedAt    time.Time   `json:"updated_at"`
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

// DefaultHome is used when nothing has been saved yet (and as the "reset" target).
func DefaultHome() HomeContent {
	return HomeContent{
		Greeting:        "Greeting",
		Headline:        "Headline",
		Description:     "Description",
		PrimaryLabel:    "PrimaryLabel",
		SecondaryLabel:  "SecondaryLabel",
		Note:            "Note",
		FloatingNote:    "FloatingNote",
		FeaturedSlug:    "FeaturedSlug",
		FeaturesLabel:   "FeaturesLabel",
		FeaturesHeading: "FeaturesHeading",
		FeaturesIntro:   "FeaturesIntro",
		Features:        []Feature{},
		AboutLabel:      "AboutLabel",
		AboutHeading:    "AboutHeading",
		AboutText:       "AboutText",
		AboutButton:     "AboutButton",
		CTAHeading:      "CTAHeading",
		CTAText:         "CTAText",
		CTAButton:       "CTAButton",
		FooterText:      "FooterText",
	}
}

func (s *Store) GetSite(ctx context.Context) (*Site, error) {
	var x Site
	var home []byte
	err := s.db.QueryRowContext(ctx, `
		SELECT name, tagline, contact_email, contact_text, primary_color,
		       logo IS NOT NULL, portrait IS NOT NULL, home, registration, updated_at
		FROM site_settings WHERE id = 1`,
	).Scan(&x.Name, &x.Tagline, &x.ContactEmail, &x.ContactText, &x.PrimaryColor,
		&x.HasLogo, &x.HasPortrait, &home, &x.Registration, &x.UpdatedAt)
	if err != nil {
		return nil, translate(err)
	}
	if len(home) > 0 {
		if err := json.Unmarshal(home, &x.Home); err != nil {
			return nil, err
		}
	}
	if x.Home.Headline == "" { // nothing saved yet (or reset) → defaults
		x.Home = DefaultHome()
	}
	return &x, nil
}

func (s *Store) UpdateSite(ctx context.Context, x *Site) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE site_settings
		SET name=$1, tagline=$2, contact_email=$3, contact_text=$4, primary_color=$5, registration=$6, updated_at=now()
		WHERE id = 1`,
		x.Name, x.Tagline, x.ContactEmail, x.ContactText, x.PrimaryColor, x.Registration)
	return err
}

// UpdateHome stores the home page content. Pass the zero value to reset to defaults.
func (s *Store) UpdateHome(ctx context.Context, h HomeContent) error {
	b, err := json.Marshal(h)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `UPDATE site_settings SET home=$1, updated_at=now() WHERE id = 1`, b)
	return err
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
