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

	if id == currentUser(r).ID {
		writeError(w, http.StatusBadGateway, "You cannot update your own user record.")
		return
	}

	u, err := h.Store.AdminUpdateUser(r.Context(), id, in.Email, in.Name, in.Verified, in.Role)
	if errors.Is(err, store.ErrDuplicate) {
		writeError(w, http.StatusConflict, "Another user already has that email.")
		return
	}

	me := currentUser(r)
	if in.Role == "owner" && me.Role != "owner" {
		writeError(w, http.StatusForbidden, "Only the owner can transfer ownership.")
		return
	}
	if id == me.ID && in.Role != me.Role {
		writeError(w, http.StatusBadRequest, "You can't change your own role.")
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

// POST /api/resumes/{slug}/unlock {password}
func (h *handlers) resumeUnlock(w http.ResponseWriter, r *http.Request) {
	res, err := h.Store.ResumeBySlug(r.Context(), r.PathValue("slug"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	var in struct {
		Password string `json:"password"`
	}
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, 400, "invalid request")
		return
	}
	if res.AccessHash == "" || bcrypt.CompareHashAndPassword([]byte(res.AccessHash), []byte(in.Password)) != nil {
		writeError(w, http.StatusUnauthorized, "Wrong password.")
		return
	}
	h.setAccessCookie(w, res)
	w.WriteHeader(http.StatusNoContent)
}

// PUT /api/me/resume/access {password}   ("" clears)
func (h *handlers) setResumeAccess(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Password string `json:"password"`
	}
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	hash := ""
	if in.Password != "" {
		if len(in.Password) < 4 {
			writeError(w, http.StatusBadRequest, "Password must be at least 4 characters.")
			return
		}
		b, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			serverError(w, r, err)
			return
		}
		hash = string(b)
	}
	if err := h.Store.SetResumeAccess(r.Context(), currentUser(r).ID, hash); err != nil {
		serverError(w, r, err)
		return
	}
	res, err := h.Store.ResumeByUser(r.Context(), currentUser(r).ID)
	h.writeItemResult(w, r, res, err)
}
