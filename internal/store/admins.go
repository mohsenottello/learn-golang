package store

import (
	"context"
	"time"
)

type Admin struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

// EnsureAdmin creates a verified admin if none exists with that email.
// Used to seed the first admin from configuration; safe to call on every start.
func (s *Store) EnsureAdmin(ctx context.Context, email, passwordHash string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO admins (email, password_hash, verified)
		VALUES ($1, $2, TRUE)
		ON CONFLICT (email) DO NOTHING`, email, passwordHash)
	return err
}

func (s *Store) AdminByEmail(ctx context.Context, email string) (*Admin, error) {
	var a Admin
	err := s.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, created_at FROM admins WHERE email = $1 AND verified`, email,
	).Scan(&a.ID, &a.Email, &a.PasswordHash, &a.CreatedAt)
	if err != nil {
		return nil, translate(err)
	}
	return &a, nil
}

func (s *Store) UpdateAdminPassword(ctx context.Context, id int64, passwordHash string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE admins SET password_hash = $2 WHERE id = $1`, id, passwordHash)
	return err
}
