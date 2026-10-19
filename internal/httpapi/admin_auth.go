package httpapi

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *handlers) adminLogin(w http.ResponseWriter, r *http.Request) {
	var in credentials
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	email := strings.ToLower(strings.TrimSpace(in.Email))

	admin, err := h.Store.AdminByEmail(r.Context(), email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(in.Password)) != nil {
		writeError(w, http.StatusUnauthorized, "Invalid email or password.")
		return
	}

	token, _ := newToken()
	expires := time.Now().Add(sessionTTL)
	if err := h.Store.CreateAdminSession(r.Context(), token, admin.ID, expires); err != nil {
		serverError(w, r, err)
		return
	}
	h.setSessionCookie(w, adminCookie, token, expires)
	writeJSON(w, http.StatusOK, admin)
}

func (h *handlers) adminLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(adminCookie); err == nil {
		h.Store.DeleteAdminSession(r.Context(), c.Value)
	}
	clearCookie(w, adminCookie)
	w.WriteHeader(http.StatusNoContent)
}

func (h *handlers) adminMe(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, currentAdmin(r))
}
