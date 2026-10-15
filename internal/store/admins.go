package store

import (
	"context"
	"time"
)

type Admin struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Verified     bool      `json:"verified"`
	CreatedAt    time.Time `json:"created_at"`
}

func (s *Store) CreateAdmin(ctx context.Context, email, passwordHash, verifyToken string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO admins (email, password_hash, verify_token) VALUES ($1, $2, $3)`,
		email, passwordHash, verifyToken)
	return translate(err)
}

// VerifyAdmin consumes a verification token. Returns ErrNotFound if it was invalid or already used.
func (s *Store) VerifyAdmin(ctx context.Context, token string) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE admins SET verified = TRUE, verify_token = NULL WHERE verify_token = $1`, token)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) AdminByEmail(ctx context.Context, email string) (*Admin, error) {
	var a Admin
	err := s.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, verified, created_at FROM admins WHERE email = $1`, email,
	).Scan(&a.ID, &a.Email, &a.PasswordHash, &a.Verified, &a.CreatedAt)
	if err != nil {
		return nil, translate(err)
	}
	return &a, nil
}
