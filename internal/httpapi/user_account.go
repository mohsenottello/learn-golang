package httpapi

// Signed-in user's own account settings.

import (
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (h *handlers) updateName(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name string `json:"name"`
	}
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	name := strings.TrimSpace(in.Name)
	if name == "" || len(name) > 100 {
		writeError(w, http.StatusBadRequest, "Name is required.")
		return
	}
	u, err := h.Store.UpdateUserName(r.Context(), currentUser(r).ID, name)
	if err != nil {
		serverError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (h *handlers) changePassword(w http.ResponseWriter, r *http.Request) {
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
	u := currentUser(r)
	full, err := h.Store.UserByEmail(r.Context(), u.Email) // need the hash; context user has none
	if err != nil || bcrypt.CompareHashAndPassword([]byte(full.PasswordHash), []byte(in.Current)) != nil {
		writeError(w, http.StatusUnauthorized, "Current password is incorrect.")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.New), bcrypt.DefaultCost)
	if err != nil {
		serverError(w, r, err)
		return
	}
	if err := h.Store.UpdatePassword(r.Context(), u.ID, string(hash)); err != nil {
		serverError(w, r, err)
		return
	}
	if c, err := r.Cookie(userCookie); err == nil {
		h.Store.DeleteOtherUserSessions(r.Context(), u.ID, c.Value)
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Password updated. Other devices have been signed out."})
}
