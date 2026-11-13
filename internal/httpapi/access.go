package httpapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/mohsenottello/golang/internal/store"
)

func accessCookieName(slug string) string { return "ra_" + slug }

func (h *handlers) accessToken(r *store.Resume) string {
	m := hmac.New(sha256.New, []byte(h.Secret))
	m.Write([]byte(r.Slug + ":" + r.AccessHash))
	return hex.EncodeToString(m.Sum(nil))
}

// unlocked reports whether this request may see the full resume.
func (h *handlers) unlocked(r *http.Request, res *store.Resume) bool {
	if res.AccessHash == "" {
		return true
	}
	c, err := r.Cookie(accessCookieName(res.Slug))
	return err == nil && hmac.Equal([]byte(c.Value), []byte(h.accessToken(res)))
}

func (h *handlers) setAccessCookie(w http.ResponseWriter, res *store.Resume) {
	http.SetCookie(w, &http.Cookie{
		Name: accessCookieName(res.Slug), Value: h.accessToken(res), Path: "/",
		HttpOnly: true, Secure: h.secureCookies(), SameSite: http.SameSiteLaxMode,
		Expires: time.Now().Add(30 * 24 * time.Hour),
	})
}
func (h *handlers) secureCookies() bool {
	return strings.HasPrefix(h.BaseURL, "https://")
}
