// Package store is the only package that talks to the database.
// Each model gets its own file (admins.go, sessions.go, users.go, ...).
package store

import (
	"context"
	"database/sql"
	"time"
)

type Site struct {
	Name         string    `json:"name"`
	Tagline      string    `json:"tagline"`
	ContactEmail string    `json:"contact_email"`
	ContactText  string    `json:"contact_text"`
	HasLogo      bool      `json:"has_logo"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (s *Store) GetSite(ctx context.Context) (*Site, error) {
	var x Site
	err := s.db.QueryRowContext(ctx, `
		SELECT name, tagline, contact_email, contact_text, logo IS NOT NULL, updated_at
		FROM site_settings WHERE id = 1`).Scan(&x.Name, &x.Tagline, &x.ContactEmail, &x.ContactText, &x.HasLogo, &x.UpdatedAt)
	return &x, translate(err)
}

func (s *Store) UpdateSite(ctx context.Context, x *Site) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE site_settings SET name=$1, tagline=$2, contact_email=$3, contact_text=$4, updated_at=now() WHERE id = 1`,
		x.Name, x.Tagline, x.ContactEmail, x.ContactText)
	return err
}

func (s *Store) SetSiteLogo(ctx context.Context, data []byte, ct string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE site_settings SET logo=$1, logo_type=$2, updated_at=now() WHERE id = 1`, data, ct)
	return err
}

func (s *Store) SiteLogo(ctx context.Context) ([]byte, string, error) {
	var data []byte
	var ct sql.NullString
	err := s.db.QueryRowContext(ctx, `SELECT logo, logo_type FROM site_settings WHERE id = 1 AND logo IS NOT NULL`).Scan(&data, &ct)
	return data, ct.String, translate(err)
}
