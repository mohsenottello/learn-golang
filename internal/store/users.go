package store

import (
	"context"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Verified     bool      `json:"verified"`
	CreatedAt    time.Time `json:"created_at"`
}

const userCols = `id, email, name, password_hash, verified, created_at`

func scanUser(row interface{ Scan(...any) error }) (*User, error) {
	var u User
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.Verified, &u.CreatedAt)
	if err != nil {
		return nil, translate(err)
	}
	return &u, nil
}

// --- registration / verification ---

func (s *Store) CreateUser(ctx context.Context, email, name, passwordHash, verifyTokenHash string) (*User, error) {
	return scanUser(s.db.QueryRowContext(ctx, `
		INSERT INTO users (email, name, password_hash, verify_token)
		VALUES ($1, $2, $3, $4)
		RETURNING `+userCols, email, name, passwordHash, verifyTokenHash))
}

func (s *Store) VerifyUser(ctx context.Context, verifyTokenHash string) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE users SET verified = TRUE, verify_token = NULL WHERE verify_token = $1`, verifyTokenHash)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// --- login ---

func (s *Store) UserByEmail(ctx context.Context, email string) (*User, error) {
	return scanUser(s.db.QueryRowContext(ctx,
		`SELECT `+userCols+` FROM users WHERE email = $1`, email))
}

// --- password reset ---

func (s *Store) SetResetToken(ctx context.Context, email, resetTokenHash string, expires time.Time) (*User, error) {
	return scanUser(s.db.QueryRowContext(ctx, `
		UPDATE users SET reset_token = $2, reset_expires_at = $3
		WHERE email = $1 AND verified
		RETURNING `+userCols, email, resetTokenHash, expires))
}

func (s *Store) UserByResetToken(ctx context.Context, resetTokenHash string) (*User, error) {
	return scanUser(s.db.QueryRowContext(ctx, `
		SELECT `+userCols+` FROM users
		WHERE reset_token = $1 AND reset_expires_at > now()`, resetTokenHash))
}

// UpdatePassword sets a new hash and invalidates any outstanding reset token.
func (s *Store) UpdatePassword(ctx context.Context, userID int64, passwordHash string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE users SET password_hash = $2, reset_token = NULL, reset_expires_at = NULL
		WHERE id = $1`, userID, passwordHash)
	return err
}

// --- admin views ---

func (s *Store) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+userCols+` FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, rows.Err()
}

func (s *Store) GetUser(ctx context.Context, id int64) (*User, error) {
	return scanUser(s.db.QueryRowContext(ctx, `SELECT `+userCols+` FROM users WHERE id = $1`, id))
}

func (s *Store) DeleteUser(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}
