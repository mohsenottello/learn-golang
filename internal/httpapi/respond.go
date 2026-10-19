package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// readJSON decodes the body. Requiring application/json doubles as a cheap
// CSRF guard: browsers cannot send that content type cross-site without a
// CORS preflight, which this API never approves.
func readJSON(w http.ResponseWriter, r *http.Request, v any) error {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return errors.New("expected application/json")
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// serverError logs the real cause and sends a generic 500 to the client.
func serverError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("500 %s %s: %v", r.Method, r.URL.Path, err)
	serverError(w, r, err)
}
