package httpapi

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/mohsenottello/golang/internal/store"
)

type ctxKey int

const (
	adminKey ctxKey = iota
	userKey
)

const (
	adminCookie = "admin_session"
	userCookie  = "user_session"
	sessionTTL  = 7 * 24 * time.Hour
)

func (h *handlers) setSessionCookie(w http.ResponseWriter, name, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   strings.HasPrefix(h.BaseURL, "https://"),
	})
}

func clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{Name: name, Value: "", Path: "/", MaxAge: -1})
}

func (h *handlers) requireAdmin(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(adminCookie)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "not logged in")
			return
		}
		admin, err := h.Store.AdminForSession(r.Context(), c.Value)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "not logged in")
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), adminKey, admin)))
	})
}

func (h *handlers) requireUser(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(userCookie)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "not logged in")
			return
		}
		user, err := h.Store.UserForSession(r.Context(), c.Value)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "not logged in")
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), userKey, user)))
	})
}

func currentAdmin(r *http.Request) *store.Admin { return r.Context().Value(adminKey).(*store.Admin) }
func currentUser(r *http.Request) *store.User   { return r.Context().Value(userKey).(*store.User) }

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(sw, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, sw.status, time.Since(start))
	})
}

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v\n%s", rec, debug.Stack())
				writeError(w, http.StatusInternalServerError, "server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
