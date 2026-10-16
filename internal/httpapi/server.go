// Package httpapi is everything that knows about HTTP: routing, middleware,
// request decoding, response encoding. It calls into store and mail; it never
// touches SQL directly.
package httpapi

import (
	"io/fs"
	"net/http"

	"github.com/mohsenottello/golang/internal/mail"
	"github.com/mohsenottello/golang/internal/store"
)

type Deps struct {
	Store   *store.Store
	Mailer  *mail.Mailer
	BaseURL string
	Static  fs.FS
}

type handlers struct{ Deps }

// NewRouter is the single map of the API.
func NewRouter(d Deps) http.Handler {
	h := &handlers{d}
	mux := http.NewServeMux()
	admin := h.requireAdmin
	user := h.requireUser

	// Users: self-service auth
	mux.HandleFunc("POST /api/auth/register", h.userRegister)
	mux.HandleFunc("POST /api/auth/verify", h.userVerify)
	mux.HandleFunc("POST /api/auth/login", h.userLogin)
	mux.HandleFunc("POST /api/auth/logout", h.userLogout)
	mux.HandleFunc("POST /api/auth/forgot-password", h.forgotPassword)
	mux.HandleFunc("POST /api/auth/reset-password", h.resetPassword)
	mux.Handle("GET /api/auth/me", user(h.userMe))

	// Admins: login only (seeded from config), plus user management
	mux.HandleFunc("POST /api/admin/login", h.adminLogin)
	mux.HandleFunc("POST /api/admin/logout", h.adminLogout)
	mux.Handle("GET /api/admin/me", admin(h.adminMe))
	mux.Handle("GET /api/admin/users", admin(h.listUsers))
	mux.Handle("GET /api/admin/users/{id}", admin(h.getUser))
	mux.Handle("DELETE /api/admin/users/{id}", admin(h.deleteUser))

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.Handle("/", spaHandler(d.Static))

	return recoverPanic(logRequests(mux))
}
