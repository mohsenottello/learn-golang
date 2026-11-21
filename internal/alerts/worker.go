// Package alerts polls prices and sends Telegram messages when a watch's threshold is crossed.
package alerts

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/mohsenottello/golang/internal/crypto"
	"github.com/mohsenottello/golang/internal/store"
	"github.com/mohsenottello/golang/internal/telegram"
)

const (
	interval = 60 * time.Second
	minGap   = 5 * time.Minute // never alert the same watch twice within this window
)

func Run(ctx context.Context, st *store.Store, tg *telegram.Client, px *crypto.Client) {
	if !tg.Enabled() {
		log.Println("alerts: TELEGRAM_BOT_TOKEN not set, worker disabled")
		return
	}
	go linkLoop(ctx, st, tg)
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			tick(ctx, st, tg, px)
		}
	}
}

func tick(ctx context.Context, st *store.Store, tg *telegram.Client, px *crypto.Client) {
	watches, err := st.AllWatchesWithChat(ctx)
	if err != nil {
		log.Println("alerts: watches:", err)
		return
	}
	if len(watches) == 0 {
		return
	}
	prices, err := px.Prices(ctx, distinct(watches))
	if err != nil {
		log.Println("alerts: prices:", err)
		return
	}
	for _, w := range watches {
		p, ok := prices[w.Symbol]
		if !ok || w.RefPrice == 0 {
			continue
		}
		change := (p - w.RefPrice) / w.RefPrice * 100
		if math.Abs(change) < w.Pct || time.Since(w.RefTime) < minGap {
			continue
		}
		arrow := "▲"
		if change < 0 {
			arrow = "▼"
		}
		msg := fmt.Sprintf("<b>%s</b> %s  %s %.2f%%\nsince %s (%s ago)",
			w.Symbol, money(p), arrow, change, money(w.RefPrice), ago(w.RefTime))
		if err := tg.Send(ctx, w.ChatID, msg); err != nil {
			log.Println("alerts: send:", err)
		}
		// Reset even if sending failed, so a Telegram hiccup can't cause a burst later.
		if err := st.ResetWatchRef(ctx, w.ID, p); err != nil {
			log.Println("alerts: reset:", err)
		}
	}
}

// linkLoop consumes "/start CODE" messages and binds the chat to the user holding that code.
func linkLoop(ctx context.Context, st *store.Store, tg *telegram.Client) {
	var offset int64
	for ctx.Err() == nil {
		updates, err := tg.Updates(ctx, offset)
		if err != nil {
			log.Println("alerts: updates:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		for _, u := range updates {
			offset = u.ID + 1
			if u.Message == nil {
				continue
			}
			f := strings.Fields(u.Message.Text)
			chat := u.Message.Chat.ID
			if len(f) == 2 && f[0] == "/start" {
				email, err := st.LinkTelegramByCode(ctx, strings.ToUpper(f[1]), chat)
				if err != nil {
					tg.Send(ctx, chat, "Unknown or expired code. Copy it again from the website.")
					continue
				}
				tg.Send(ctx, chat, "Linked to <b>"+email+"</b> ✔ You'll receive price alerts here.")
			} else {
				tg.Send(ctx, chat, "Manage your alerts on the website; send /start CODE from there to link this chat.")
			}
		}
	}
}

func distinct(ws []store.Watch) []string {
	seen := map[string]bool{}
	var out []string
	for _, w := range ws {
		if !seen[w.Symbol] {
			seen[w.Symbol] = true
			out = append(out, w.Symbol)
		}
	}
	return out
}

func money(p float64) string {
	switch {
	case p >= 1000:
		return fmt.Sprintf("$%.0f", p)
	case p >= 1:
		return fmt.Sprintf("$%.2f", p)
	default:
		return fmt.Sprintf("$%.4g", p)
	}
}

func ago(t time.Time) string { return time.Since(t).Round(time.Minute).String() }
