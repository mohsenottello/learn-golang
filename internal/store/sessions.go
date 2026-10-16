package store

import (
	"context"
	"time"
)

// Admin sessions.

func (s *Store) CreateAdminSession(ctx context.Context, token string, adminID int64, expires time.Time) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO sessions (token, admin_id, expires_at) VALUES ($1, $2, $3)`,
		token, adminID, expires)
	return err
}

func (s *Store) AdminForSession(ctx context.Context, token string) (*Admin, error) {
	var a Admin
	err := s.db.QueryRowContext(ctx, `
		SELECT a.id, a.email, a.created_at
		FROM sessions s JOIN admins a ON a.id = s.admin_id
		WHERE s.token = $1 AND s.expires_at > now()`, token,
	).Scan(&a.ID, &a.Email, &a.CreatedAt)
	if err != nil {
		return nil, translate(err)
	}
	return &a, nil
}

func (s *Store) DeleteAdminSession(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM sessions WHERE token = $1`, token)
	return err
}

// User sessions.

func (s *Store) CreateUserSession(ctx context.Context, token string, userID int64, expires time.Time) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO user_sessions (token, user_id, expires_at) VALUES ($1, $2, $3)`,
		token, userID, expires)
	return err
}

func (s *Store) UserForSession(ctx context.Context, token string) (*User, error) {
	var u User
	err := s.db.QueryRowContext(ctx, `
		SELECT u.id, u.email, u.name, u.verified, u.created_at
		FROM user_sessions s JOIN users u ON u.id = s.user_id
		WHERE s.token = $1 AND s.expires_at > now()`, token,
	).Scan(&u.ID, &u.Email, &u.Name, &u.Verified, &u.CreatedAt)
	if err != nil {
		return nil, translate(err)
	}
	return &u, nil
}

func (s *Store) DeleteUserSession(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM user_sessions WHERE token = $1`, token)
	return err
}

// DeleteAllUserSessions logs a user out everywhere (used after a password reset).
func (s *Store) DeleteAllUserSessions(ctx context.Context, userID int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM user_sessions WHERE user_id = $1`, userID)
	return err
}

// DeleteOtherUserSessions logs the user out everywhere except the given session.
func (s *Store) DeleteOtherUserSessions(ctx context.Context, userID int64, keepToken string) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM user_sessions WHERE user_id = $1 AND token <> $2`, userID, keepToken)
	return err
}
