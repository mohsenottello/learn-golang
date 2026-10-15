package httpapi

// Example resource. The pattern for every new API:
//   decode input → validate → call store → map store errors → encode output.

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/mohsenottello/golang/internal/store"
)

func (h *handlers) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Store.ListUsers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server error")
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *handlers) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	u, err := h.Store.GetUser(r.Context(), id)
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, "user not found")
	case err != nil:
		writeError(w, http.StatusInternalServerError, "server error")
	default:
		writeJSON(w, http.StatusOK, u)
	}
}

func (h *handlers) createUser(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	in.Name = strings.TrimSpace(in.Name)
	if !strings.Contains(in.Email, "@") || in.Name == "" {
		writeError(w, http.StatusBadRequest, "email and name are required")
		return
	}

	u, err := h.Store.CreateUser(r.Context(), in.Email, in.Name)
	switch {
	case errors.Is(err, store.ErrDuplicate):
		writeError(w, http.StatusConflict, "a user with that email already exists")
	case err != nil:
		writeError(w, http.StatusInternalServerError, "server error")
	default:
		writeJSON(w, http.StatusCreated, u)
	}
}

func (h *handlers) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	err = h.Store.DeleteUser(r.Context(), id)
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, "user not found")
	case err != nil:
		writeError(w, http.StatusInternalServerError, "server error")
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
