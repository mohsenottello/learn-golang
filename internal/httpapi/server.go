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

// Deps is everything the handlers need. Passed in from main, never global.
type Deps struct {
	Store   *store.Store
	Mailer  *mail.Mailer
	BaseURL string
	Static  fs.FS // the built React app
}

type handlers struct{ Deps }

// NewRouter is the single map of the API. Add new routes here.
func NewRouter(d Deps) http.Handler {
	h := &handlers{d}
	mux := http.NewServeMux()
	auth := h.requireAdmin // protected routes wrap with this

	// Auth
	mux.HandleFunc("POST /api/admin/register", h.register)
	mux.HandleFunc("POST /api/admin/verify", h.verify)
	mux.HandleFunc("POST /api/admin/login", h.login)
	mux.HandleFunc("POST /api/admin/logout", h.logout)
	mux.Handle("GET /api/admin/me", auth(h.me))

	// Users (example resource)
	mux.Handle("GET /api/users", auth(h.listUsers))
	mux.Handle("POST /api/users", auth(h.createUser))
	mux.Handle("GET /api/users/{id}", auth(h.getUser))
	mux.Handle("DELETE /api/users/{id}", auth(h.deleteUser))

	// Health check for Docker / load balancers
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// React app for everything else
	mux.Handle("/", spaHandler(d.Static))

	return recoverPanic(logRequests(mux))
}
