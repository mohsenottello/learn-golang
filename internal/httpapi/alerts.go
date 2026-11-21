package httpapi

import (
	"crypto/rand"
	"encoding/base32"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var symbolRe = regexp.MustCompile(`^[A-Z0-9]{2,10}$`)

func newLinkCode() string {
	b := make([]byte, 5)
	rand.Read(b)
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b) // 8 chars, e.g. "K7Q2MZ4A"
}

// GET /api/me/alerts
func (h *handlers) getAlerts(w http.ResponseWriter, r *http.Request) {
	uid := currentUser(r).ID
	linked, code, err := h.Store.TelegramStatus(r.Context(), uid)
	if err != nil {
		serverError(w, r, err)
		return
	}
	if !linked && code == "" {
		code = newLinkCode()
		if err := h.Store.SetTelegramLinkCode(r.Context(), uid, code); err != nil {
			serverError(w, r, err)
			return
		}
	}
	watches, err := h.Store.ListWatchesByUser(r.Context(), uid)
	if err != nil {
		serverError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"enabled":         h.Telegram.Enabled(),
		"telegram_linked": linked,
		"link_code":       code,
		"bot_username":    h.BotUsername,
		"watches":         watches,
	})
}

// POST /api/me/alerts/watches {symbol, pct}
func (h *handlers) addWatch(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Symbol string  `json:"symbol"`
		Pct    float64 `json:"pct"`
	}
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	sym := strings.ToUpper(strings.TrimSpace(in.Symbol))
	if !symbolRe.MatchString(sym) {
		writeError(w, http.StatusBadRequest, "Symbol must be letters/digits, e.g. BTC.")
		return
	}
	if in.Pct < 0.1 || in.Pct > 100 {
		writeError(w, http.StatusBadRequest, "Percent must be between 0.1 and 100.")
		return
	}
	price, err := h.Prices.Price(r.Context(), sym)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Unknown symbol on Binance: "+sym)
		return
	}
	wt, err := h.Store.UpsertWatch(r.Context(), currentUser(r).ID, sym, in.Pct, price)
	h.writeItemResult(w, r, wt, err)
}

// DELETE /api/me/alerts/watches/{id}
func (h *handlers) deleteWatch(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := h.Store.DeleteWatch(r.Context(), currentUser(r).ID, id); err != nil {
		h.writeItemResult(w, r, nil, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/me/alerts/telegram
func (h *handlers) unlinkTelegram(w http.ResponseWriter, r *http.Request) {
	if err := h.Store.UnlinkTelegram(r.Context(), currentUser(r).ID); err != nil {
		serverError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
