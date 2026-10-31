package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/mohsenottello/golang/internal/translate"
)

// POST /api/translate  {texts: [...], source: "de", target: "en"}  → {texts: [...]}
func (h *handlers) translateTexts(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Texts  []string `json:"texts"`
		Source string   `json:"source"`
		Target string   `json:"target"`
	}
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if !validLangs[in.Source] || !validLangs[in.Target] || in.Source == in.Target {
		writeError(w, http.StatusBadRequest, "source and target must be de/en and differ")
		return
	}
	total := 0
	for i := range in.Texts {
		in.Texts[i] = strings.TrimSpace(in.Texts[i])
		total += len(in.Texts[i])
	}
	if len(in.Texts) == 0 || len(in.Texts) > 50 || total > 20_000 {
		writeError(w, http.StatusBadRequest, "up to 50 texts / 20,000 characters per request")
		return
	}
	out, err := h.Translator.Translate(r.Context(), in.Texts, in.Source, in.Target)
	switch {
	case errors.Is(err, translate.ErrDisabled):
		writeError(w, http.StatusNotImplemented, "Translation is not configured on this server.")
	case err != nil:
		serverError(w, r, err)
	default:
		writeJSON(w, http.StatusOK, map[string]any{"texts": out})
	}
}
