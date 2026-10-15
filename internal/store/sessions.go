package store

import (
	"context"
	"time"
)

func (s *Store) CreateSession(ctx context.Context, token string, adminID int64, expires time.Time) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO sessions (token, admin_id, expires_at) VALUES ($1, $2, $3)`,
		token, adminID, expires)
	return err
}

// AdminForSession returns the admin for a live session token, or ErrNotFound.
func (s *Store) AdminForSession(ctx context.Context, token string) (*Admin, error) {
	var a Admin
	err := s.db.QueryRowContext(ctx, `
		SELECT a.id, a.email, a.verified, a.created_at
		FROM sessions s JOIN admins a ON a.id = s.admin_id
		WHERE s.token = $1 AND s.expires_at > now()`, token,
	).Scan(&a.ID, &a.Email, &a.Verified, &a.CreatedAt)
	if err != nil {
		return nil, translate(err)
	}
	return &a, nil
}

func (s *Store) DeleteSession(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM sessions WHERE token = $1`, token)
	return err
}
