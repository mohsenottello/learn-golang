package httpapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mohsenottello/golang/internal/store"
)

const accessTTL = 5 * time.Minute

func accessCookieName(slug string) string { return "ra_" + slug }

func (h *handlers) sign(res *store.Resume, ts string) string {
	m := hmac.New(sha256.New, []byte(h.Secret))
	m.Write([]byte(res.Slug + ":" + res.AccessHash + ":" + ts))
	return hex.EncodeToString(m.Sum(nil))
}

// cookie value: "<unix-seconds>.<hmac>"
func (h *handlers) accessToken(res *store.Resume) string {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	return ts + "." + h.sign(res, ts)
}

func (h *handlers) unlocked(r *http.Request, res *store.Resume) bool {
	if res.AccessHash == "" {
		return true
	}
	c, err := r.Cookie(accessCookieName(res.Slug))
	if err != nil {
		return false
	}
	ts, sig, ok := strings.Cut(c.Value, ".")
	if !ok {
		return false
	}
	issued, err := strconv.ParseInt(ts, 10, 64)
	if err != nil || time.Since(time.Unix(issued, 0)) > accessTTL {
		return false
	}
	return hmac.Equal([]byte(sig), []byte(h.sign(res, ts)))
}

func (h *handlers) setAccessCookie(w http.ResponseWriter, res *store.Resume) {
	http.SetCookie(w, &http.Cookie{
		Name: accessCookieName(res.Slug), Value: h.accessToken(res), Path: "/",
		HttpOnly: true, Secure: h.secureCookies(), SameSite: http.SameSiteLaxMode,
		MaxAge: int(accessTTL.Seconds()),
	})
}

func (h *handlers) secureCookies() bool {
	return strings.HasPrefix(h.BaseURL, "https://")
}
