package httpapi

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/mohsenottello/golang/internal/store"
)

const resetTTL = time.Hour

// Same response whether or not the email exists, so accounts can't be enumerated.
const neutralRegisterMsg = "If that address is new, we've sent you a verification link."
const neutralResetMsg = "If an account exists for that address, we've sent a reset link."

func validEmail(s string) bool { return strings.Contains(s, "@") && len(s) <= 254 }

func (h *handlers) userRegister(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	email := strings.ToLower(strings.TrimSpace(in.Email))
	name := strings.TrimSpace(in.Name)
	if !validEmail(email) || name == "" || len(in.Password) < 8 {
		writeError(w, http.StatusBadRequest, "Name, a valid email and a password of at least 8 characters are required.")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		serverError(w, r, err)
		return
	}

	raw, tokenHash := newToken()
	_, err = h.Store.CreateUser(r.Context(), email, name, string(hash), tokenHash)
	switch {
	case errors.Is(err, store.ErrDuplicate):
		// say nothing
	case err != nil:
		log.Println("register:", err)
		serverError(w, r, err)
		return
	default:
		link := h.BaseURL + "/verify?token=" + raw
		go h.Mailer.Send(email, "Verify your account",
			"Hi "+name+",\n\nClick the link to verify your account:\n\n"+link+"\n")
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": neutralRegisterMsg})
}

func (h *handlers) userVerify(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Token string `json:"token"`
	}
	if err := readJSON(w, r, &in); err != nil || in.Token == "" {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	err := h.Store.VerifyUser(r.Context(), hashToken(in.Token))
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusBadRequest, "This verification link is invalid or was already used.")
	case err != nil:
		serverError(w, r, err)
	default:
		writeJSON(w, http.StatusOK, map[string]string{"message": "Your account is verified. You can log in now."})
	}
}

func (h *handlers) userLogin(w http.ResponseWriter, r *http.Request) {
	var in credentials
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	email := strings.ToLower(strings.TrimSpace(in.Email))

	user, err := h.Store.UserByEmail(r.Context(), email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
		writeError(w, http.StatusUnauthorized, "Invalid email or password.")
		return
	}
	if !user.Verified {
		writeError(w, http.StatusForbidden, "Please verify your email first.")
		return
	}

	token, _ := newToken()
	expires := time.Now().Add(sessionTTL)
	if err := h.Store.CreateUserSession(r.Context(), token, user.ID, expires); err != nil {
		serverError(w, r, err)
		return
	}
	h.setSessionCookie(w, userCookie, token, expires)
	writeJSON(w, http.StatusOK, user)
}

func (h *handlers) userLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(userCookie); err == nil {
		h.Store.DeleteUserSession(r.Context(), c.Value)
	}
	clearCookie(w, userCookie)
	w.WriteHeader(http.StatusNoContent)
}

func (h *handlers) userMe(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, currentUser(r))
}

// --- password recovery ---

func (h *handlers) forgotPassword(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email string `json:"email"`
	}
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	email := strings.ToLower(strings.TrimSpace(in.Email))

	raw, tokenHash := newToken()
	user, err := h.Store.SetResetToken(r.Context(), email, tokenHash, time.Now().Add(resetTTL))
	switch {
	case errors.Is(err, store.ErrNotFound):
		// unknown or unverified email: say nothing
	case err != nil:
		log.Println("forgot password:", err)
		serverError(w, r, err)
		return
	default:
		link := h.BaseURL + "/reset-password?token=" + raw
		go h.Mailer.Send(user.Email, "Reset your password",
			"Hi "+user.Name+",\n\nClick the link to choose a new password. It is valid for one hour:\n\n"+link+
				"\n\nIf you didn't request this, you can ignore this email.\n")
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": neutralResetMsg})
}

func (h *handlers) resetPassword(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := readJSON(w, r, &in); err != nil || in.Token == "" {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if len(in.Password) < 8 {
		writeError(w, http.StatusBadRequest, "Password must be at least 8 characters.")
		return
	}

	user, err := h.Store.UserByResetToken(r.Context(), hashToken(in.Token))
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusBadRequest, "This reset link is invalid or has expired.")
		return
	case err != nil:
		serverError(w, r, err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		serverError(w, r, err)
		return
	}
	if err := h.Store.UpdatePassword(r.Context(), user.ID, string(hash)); err != nil {
		serverError(w, r, err)
		return
	}
	// Log out every existing session: whoever had the old password is gone.
	h.Store.DeleteAllUserSessions(r.Context(), user.ID)

	writeJSON(w, http.StatusOK, map[string]string{"message": "Password updated. You can log in now."})
}
