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
		INSERT INTO users (email, password_hash, verified, role, name)
		VALUES ($1, $2, TRUE, 'owner', $3)
		ON CONFLICT (email) DO NOTHING`, email, passwordHash, "Admin")
	return err
}

func (s *Store) AdminByEmail(ctx context.Context, email string) (*Admin, error) {
	var a Admin
	err := s.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, created_at FROM users WHERE email = $1 AND (role = 'admin' or role = 'owner')`, email,
	).Scan(&a.ID, &a.Email, &a.PasswordHash, &a.CreatedAt)
	if err != nil {
		return nil, translate(err)
	}
	return &a, nil
}
