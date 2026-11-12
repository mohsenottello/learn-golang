package httpapi

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/mohsenottello/golang/internal/store"
)

// GET /api/resumes/{slug}/pdf?lang=de|en
func (h *handlers) resumePDF(w http.ResponseWriter, r *http.Request) {
	if h.GotenbergURL == "" {
		writeError(w, http.StatusNotImplemented, "PDF export is not configured on this server.")
		return
	}
	slug := r.PathValue("slug")
	res, err := h.Store.ResumeBySlug(r.Context(), slug) // published only
	if errors.Is(err, store.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		serverError(w, r, err)
		return
	}
	lang := r.URL.Query().Get("lang")
	if !validLangs[lang] {
		lang = res.DefaultLang
	}

	// Chrome inside Gotenberg fetches the page from the app container directly.
	pageURL := fmt.Sprintf("%s/%s/%s?print=1", h.InternalURL, slug, lang)

	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	for k, v := range map[string]string{
		"url":         pageURL,
		"paperWidth":  "8.27", // A4, inches
		"paperHeight": "11.69",
		"marginTop":   "0", "marginBottom": "0", "marginLeft": "0", "marginRight": "0",
		"printBackground":   "true",
		"preferCssPageSize": "true",
		"waitForExpression": `document.body.dataset.ready === "1"`,
	} {
		mw.WriteField(k, v)
	}
	mw.Close()

	req, _ := http.NewRequestWithContext(r.Context(), "POST", h.GotenbergURL+"/forms/chromium/convert/url", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		serverError(w, r, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		serverError(w, r, fmt.Errorf("gotenberg %d: %s", resp.StatusCode, msg))
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s_CV_%s.pdf"`, slug, lang))
	w.Header().Set("Cache-Control", "private, max-age=0")
	io.Copy(w, resp.Body)
}
