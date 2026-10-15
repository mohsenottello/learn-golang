package store

// Example model to show the pattern. Copy this file for the next one.

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Store) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, email, name, created_at FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{} // not nil, so JSON is [] rather than null when empty
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (s *Store) GetUser(ctx context.Context, id int64) (*User, error) {
	var u User
	err := s.db.QueryRowContext(ctx,
		`SELECT id, email, name, created_at FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, translate(err)
	}
	return &u, nil
}

func (s *Store) CreateUser(ctx context.Context, email, name string) (*User, error) {
	var u User
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id, email, name, created_at`,
		email, name,
	).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, translate(err)
	}
	return &u, nil
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
