package store

import (
	"context"
	"time"
)

type Watch struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Symbol    string    `json:"symbol"`
	Pct       float64   `json:"pct"`
	RefPrice  float64   `json:"ref_price"`
	RefTime   time.Time `json:"ref_time"`
	CreatedAt time.Time `json:"created_at"`
	ChatID    int64     `json:"-"` // filled by AllWatchesWithChat only
}

const watchCols = `id, user_id, symbol, pct, ref_price, ref_time, created_at`

func scanWatch(row interface{ Scan(...any) error }) (*Watch, error) {
	var w Watch
	err := row.Scan(&w.ID, &w.UserID, &w.Symbol, &w.Pct, &w.RefPrice, &w.RefTime, &w.CreatedAt)
	if err != nil {
		return nil, translate(err)
	}
	return &w, nil
}

func (s *Store) ListWatchesByUser(ctx context.Context, userID int64) ([]Watch, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+watchCols+` FROM crypto_watches WHERE user_id=$1 ORDER BY symbol`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Watch{}
	for rows.Next() {
		w, err := scanWatch(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *w)
	}
	return out, rows.Err()
}

// UpsertWatch creates or replaces the user's watch for a symbol, resetting its reference price.
func (s *Store) UpsertWatch(ctx context.Context, userID int64, symbol string, pct, refPrice float64) (*Watch, error) {
	return scanWatch(s.db.QueryRowContext(ctx, `
		INSERT INTO crypto_watches (user_id, symbol, pct, ref_price, ref_time)
		VALUES ($1, $2, $3, $4, now())
		ON CONFLICT (user_id, symbol) DO UPDATE SET pct = EXCLUDED.pct, ref_price = EXCLUDED.ref_price, ref_time = now()
		RETURNING `+watchCols, userID, symbol, pct, refPrice))
}

func (s *Store) DeleteWatch(ctx context.Context, userID, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM crypto_watches WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	return rowsAffectedOrNotFound(res)
}

// AllWatchesWithChat returns every watch whose owner has linked Telegram — the worker's input.
func (s *Store) AllWatchesWithChat(ctx context.Context) ([]Watch, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT w.id, w.user_id, w.symbol, w.pct, w.ref_price, w.ref_time, w.created_at, u.telegram_chat_id
		FROM crypto_watches w JOIN users u ON u.id = w.user_id
		WHERE u.telegram_chat_id IS NOT NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Watch{}
	for rows.Next() {
		var w Watch
		if err := rows.Scan(&w.ID, &w.UserID, &w.Symbol, &w.Pct, &w.RefPrice, &w.RefTime, &w.CreatedAt, &w.ChatID); err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

// ResetWatchRef moves the reference to the current price after an alert.
func (s *Store) ResetWatchRef(ctx context.Context, id int64, price float64) error {
	_, err := s.db.ExecContext(ctx, `UPDATE crypto_watches SET ref_price=$2, ref_time=now() WHERE id=$1`, id, price)
	return err
}

// ---- Telegram linking (users table) ----

func (s *Store) SetTelegramLinkCode(ctx context.Context, userID int64, code string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE users SET telegram_link_code=$2 WHERE id=$1`, userID, code)
	return err
}

// LinkTelegramByCode binds a chat to the user holding this code and consumes the code.
// Returns the user's email for the confirmation message.
func (s *Store) LinkTelegramByCode(ctx context.Context, code string, chatID int64) (string, error) {
	var email string
	err := s.db.QueryRowContext(ctx, `
		UPDATE users SET telegram_chat_id=$2, telegram_link_code=NULL
		WHERE telegram_link_code=$1 RETURNING email`, code, chatID).Scan(&email)
	return email, translate(err)
}

func (s *Store) TelegramStatus(ctx context.Context, userID int64) (linked bool, code string, err error) {
	var chat *int64
	var c *string
	err = s.db.QueryRowContext(ctx, `SELECT telegram_chat_id, telegram_link_code FROM users WHERE id=$1`, userID).Scan(&chat, &c)
	if err != nil {
		return false, "", translate(err)
	}
	if c != nil {
		code = *c
	}
	return chat != nil, code, nil
}

func (s *Store) UnlinkTelegram(ctx context.Context, userID int64) error {
	_, err := s.db.ExecContext(ctx, `UPDATE users SET telegram_chat_id=NULL WHERE id=$1`, userID)
	return err
}
