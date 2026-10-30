package httpapi

// Admin area: user management, resume moderation, admin password.

import (
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/mohsenottello/golang/internal/store"
)

// GET /api/admin/users?q=  — list (optionally filtered) users
func (h *handlers) adminListUsers(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	users, err := h.Store.SearchUsers(r.Context(), q)
	h.writeItemResult(w, r, users, err)
}

// PUT /api/admin/users/{id}  {email, name, verified}
func (h *handlers) adminUpdateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	var in struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Role     string `json:"role"`
		Verified bool   `json:"verified"`
	}
	if err := readJSON(w, r, &in); !ok || err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	in.Name = strings.TrimSpace(in.Name)
	if !validEmail(in.Email) || in.Name == "" {
		writeError(w, http.StatusBadRequest, "A valid email and a name are required.")
		return
	}
	in.Role = strings.ToLower(strings.TrimSpace(in.Role))
	if in.Role != "user" && in.Role != "admin" && in.Role != "owner" {
		writeError(w, http.StatusBadRequest, "Role must be one of: user, admin, owner.")
		return
	}
	u, err := h.Store.AdminUpdateUser(r.Context(), id, in.Email, in.Name, in.Verified, in.Role)
	if errors.Is(err, store.ErrDuplicate) {
		writeError(w, http.StatusConflict, "Another user already has that email.")
		return
	}
	h.writeItemResult(w, r, u, err)
}

// GET /api/admin/resumes
func (h *handlers) adminListResumes(w http.ResponseWriter, r *http.Request) {
	list, err := h.Store.ListResumesForAdmin(r.Context())
	h.writeItemResult(w, r, list, err)
}

// PUT /api/admin/resumes/{id}/published  {published}   ({id} is the owner's user id)
func (h *handlers) adminSetPublished(w http.ResponseWriter, r *http.Request) {
	userID, ok := pathID(r, "id")
	var in struct {
		Published bool `json:"published"`
	}
	if err := readJSON(w, r, &in); !ok || err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	h.writeItemResult(w, r, nil, h.Store.SetResumePublished(r.Context(), userID, in.Published))
}

// PUT /api/admin/users/{id}/password  {new_password}
func (h *handlers) adminSetUserPassword(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	var in struct {
		New string `json:"new_password"`
	}
	if err := readJSON(w, r, &in); !ok || err != nil || len(in.New) < 8 {
		writeError(w, http.StatusBadRequest, "Password must be at least 8 characters.")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.New), bcrypt.DefaultCost)
	if err != nil {
		serverError(w, r, err)
		return
	}
	if err := h.Store.UpdatePassword(r.Context(), id, string(hash)); err != nil {
		serverError(w, r, err)
		return
	}
	h.Store.DeleteAllUserSessions(r.Context(), id) // they must log in again everywhere
	w.WriteHeader(http.StatusNoContent)
}
