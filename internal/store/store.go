// Package store is the only package that talks to the database.
// Each model gets its own file (admins.go, sessions.go, users.go, ...).
package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("already exists")
)

type Store struct {
	db *sql.DB
}

// Open connects and waits (up to ~30s) for the database to accept connections,
// which matters in Docker where Postgres may still be starting.
func Open(ctx context.Context, url string) (*Store, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	for i := 0; i < 15; i++ {
		if err = db.PingContext(ctx); err == nil {
			return &Store{db: db}, nil
		}
		log.Println("waiting for database:", err)
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("database not reachable: %w", err)
}

func (s *Store) Close() error { return s.db.Close() }

// Migrate applies any pending SQL migrations from the embedded migrations/ folder.
func (s *Store) Migrate() error {
	goose.SetBaseFS(migrationsFS)
	goose.SetLogger(goose.NopLogger())
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(s.db, "migrations")
}

// translate maps driver errors onto the store's sentinel errors so handlers
// never need to know about Postgres error codes.
func translate(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
		return ErrDuplicate
	}
	return err
}
