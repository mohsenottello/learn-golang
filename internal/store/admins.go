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

// EnsureOwner creates the bootstrap owner account unless an owner already exists.
// It never modifies an existing owner, so the ADMIN_* env vars stop mattering after first boot.
func (s *Store) EnsureOwner(ctx context.Context, email, name, passwordHash string) (created bool, err error) {
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO users (email, name, password_hash, verified, role)
		SELECT $1, $2, $3, TRUE, 'owner'
		WHERE NOT EXISTS (SELECT 1 FROM users WHERE role = 'owner')
		ON CONFLICT (email) DO NOTHING`, email, name, passwordHash)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n == 1, nil
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
