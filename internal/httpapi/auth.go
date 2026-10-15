package httpapi

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/mohsenottello/golang/internal/store"
)

const sessionTTL = 7 * 24 * time.Hour

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *handlers) register(w http.ResponseWriter, r *http.Request) {
	var in credentials
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	email := strings.ToLower(strings.TrimSpace(in.Email))
	if !strings.Contains(email, "@") || len(in.Password) < 8 {
		writeError(w, http.StatusBadRequest, "Enter a valid email and a password of at least 8 characters.")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server error")
		return
	}

	token := randomToken()
	err = h.Store.CreateAdmin(r.Context(), email, string(hash), token)
	switch {
	case errors.Is(err, store.ErrDuplicate):
		// Say nothing: don't reveal which emails exist.
	case err != nil:
		log.Println("register:", err)
		writeError(w, http.StatusInternalServerError, "server error")
		return
	default:
		link := h.BaseURL + "/admin/verify?token=" + token
		go h.Mailer.Send(email, "Verify your admin account",
			"Click the link to verify your account:\n\n"+link+"\n")
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "If that address is new, we've sent you a verification link.",
	})
}

func (h *handlers) verify(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Token string `json:"token"`
	}
	if err := readJSON(w, r, &in); err != nil || in.Token == "" {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	err := h.Store.VerifyAdmin(r.Context(), in.Token)
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusBadRequest, "This verification link is invalid or was already used.")
	case err != nil:
		writeError(w, http.StatusInternalServerError, "server error")
	default:
		writeJSON(w, http.StatusOK, map[string]string{"message": "Your account is verified."})
	}
}

func (h *handlers) login(w http.ResponseWriter, r *http.Request) {
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
	if !admin.Verified {
		writeError(w, http.StatusForbidden, "Please verify your email first.")
		return
	}

	token := randomToken()
	expires := time.Now().Add(sessionTTL)
	if err := h.Store.CreateSession(r.Context(), token, admin.ID, expires); err != nil {
		writeError(w, http.StatusInternalServerError, "server error")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   strings.HasPrefix(h.BaseURL, "https://"),
	})
	writeJSON(w, http.StatusOK, admin)
}

func (h *handlers) logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("session"); err == nil {
		h.Store.DeleteSession(r.Context(), c.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: "session", Value: "", Path: "/", MaxAge: -1})
	w.WriteHeader(http.StatusNoContent)
}

func (h *handlers) me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, currentAdmin(r))
}

func randomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
