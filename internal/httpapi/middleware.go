package httpapi

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/mohsenottello/golang/internal/store"
)

type ctxKey int

const adminKey ctxKey = iota

// requireAdmin rejects requests without a valid session and puts the admin
// into the request context for the wrapped handler.
func (h *handlers) requireAdmin(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
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

func currentAdmin(r *http.Request) *store.Admin {
	return r.Context().Value(adminKey).(*store.Admin)
}

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
