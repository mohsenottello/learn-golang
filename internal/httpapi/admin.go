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
	u, err := h.Store.AdminUpdateUser(r.Context(), id, in.Email, in.Name, in.Verified)
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

// PUT /api/admin/password  {current_password, new_password}
func (h *handlers) adminChangePassword(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Current string `json:"current_password"`
		New     string `json:"new_password"`
	}
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if len(in.New) < 8 {
		writeError(w, http.StatusBadRequest, "New password must be at least 8 characters.")
		return
	}
	me := currentAdmin(r)
	full, err := h.Store.AdminByEmail(r.Context(), me.Email) // context admin carries no hash
	if err != nil || bcrypt.CompareHashAndPassword([]byte(full.PasswordHash), []byte(in.Current)) != nil {
		writeError(w, http.StatusUnauthorized, "Current password is incorrect.")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.New), bcrypt.DefaultCost)
	if err != nil {
		serverError(w, r, err)
		return
	}
	if err := h.Store.UpdateAdminPassword(r.Context(), me.ID, string(hash)); err != nil {
		serverError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Password updated."})
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
