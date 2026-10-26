package httpapi

// Site-wide settings: public read (used by the nav and contact page),
// admin write. One row, see store/site.go.

import (
	"io"
	"net/http"
	"strings"

	"github.com/mohsenottello/golang/internal/store"
)

// GET /api/site — public
func (h *handlers) getSite(w http.ResponseWriter, r *http.Request) {
	site, err := h.Store.GetSite(r.Context())
	h.writeItemResult(w, r, site, err)
}

// GET /api/site/logo — public
func (h *handlers) siteLogo(w http.ResponseWriter, r *http.Request) {
	data, ct, err := h.Store.SiteLogo(r.Context())
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", ct)
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.Write(data)
}

// PUT /api/admin/site — {name, tagline, contact_email, contact_text}
func (h *handlers) adminUpdateSite(w http.ResponseWriter, r *http.Request) {
	var in store.Site
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Tagline = strings.TrimSpace(in.Tagline)
	in.ContactEmail = strings.ToLower(strings.TrimSpace(in.ContactEmail))
	in.ContactText = strings.TrimSpace(in.ContactText)

	if in.Name == "" || len(in.Name) > 60 {
		writeError(w, http.StatusBadRequest, "Site name is required (max 60 characters).")
		return
	}
	if in.ContactEmail != "" && !validEmail(in.ContactEmail) {
		writeError(w, http.StatusBadRequest, "Contact email is not valid.")
		return
	}
	if err := h.Store.UpdateSite(r.Context(), &in); err != nil {
		serverError(w, r, err)
		return
	}
	site, err := h.Store.GetSite(r.Context()) // return the saved state, incl. has_logo/updated_at
	h.writeItemResult(w, r, site, err)
}

const maxLogoBytes = 512 << 10 // 512 KB

// PUT /api/admin/site/logo — multipart field "logo"; PNG, JPEG, WebP or SVG
func (h *handlers) adminUploadLogo(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxLogoBytes+4096)
	if err := r.ParseMultipartForm(maxLogoBytes); err != nil {
		writeError(w, http.StatusBadRequest, "Logo must be smaller than 512 KB.")
		return
	}
	f, _, err := r.FormFile("logo")
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing 'logo' field")
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not read file")
		return
	}

	ct := http.DetectContentType(data)
	if strings.HasPrefix(ct, "text/") && strings.Contains(strings.ToLower(string(data[:min(len(data), 512)])), "<svg") {
		ct = "image/svg+xml" // DetectContentType doesn't recognise SVG
	}
	switch ct {
	case "image/png", "image/jpeg", "image/webp", "image/svg+xml":
	default:
		writeError(w, http.StatusBadRequest, "Only PNG, JPEG, WebP or SVG logos are allowed.")
		return
	}
	if err := h.Store.SetSiteLogo(r.Context(), data, ct); err != nil {
		serverError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
