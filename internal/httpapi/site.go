package httpapi

// Site-wide settings: public read (used by the nav, home and contact page),
// admin write. One row, see store/site.go.

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/mohsenottello/golang/internal/store"
)

var hexColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
var featureColors = map[string]bool{"teal": true, "purple": true, "pink": true, "blue": true}

// GET /api/site — public
func (h *handlers) getSite(w http.ResponseWriter, r *http.Request) {
	site, err := h.Store.GetSite(r.Context())
	h.writeItemResult(w, r, site, err)
}

// GET /api/admin/site — the full settings for the admin panel (same shape for now).
func (h *handlers) adminGetSite(w http.ResponseWriter, r *http.Request) {
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

// GET /api/site/portrait — public
func (h *handlers) sitePortrait(w http.ResponseWriter, r *http.Request) {
	data, ct, err := h.Store.Portrait(r.Context())
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", ct)
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.Write(data)
}

// PUT /api/admin/site
func (h *handlers) adminUpdateSite(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name          string `json:"name"`
		Tagline       string `json:"tagline"`
		TaglineEn     string `json:"tagline_en"`
		ContactEmail  string `json:"contact_email"`
		ContactText   string `json:"contact_text"`
		ContactTextEn string `json:"contact_text_en"`
		PrimaryColor  string `json:"primary_color"`
		Registration  bool   `json:"registration"`
		DefaultLang   string `json:"default_lang"`
	}
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	s := store.Site{
		Name:          strings.TrimSpace(in.Name),
		Tagline:       strings.TrimSpace(in.Tagline),
		TaglineEn:     strings.TrimSpace(in.TaglineEn),
		ContactEmail:  strings.ToLower(strings.TrimSpace(in.ContactEmail)),
		ContactText:   strings.TrimSpace(in.ContactText),
		ContactTextEn: strings.TrimSpace(in.ContactTextEn),
		PrimaryColor:  strings.TrimSpace(in.PrimaryColor),
		Registration:  in.Registration,
		DefaultLang:   in.DefaultLang,
	}
	if s.Name == "" || len(s.Name) > 60 {
		writeError(w, http.StatusBadRequest, "Site name is required (max 60 characters).")
		return
	}
	if s.ContactEmail != "" && !validEmail(s.ContactEmail) {
		writeError(w, http.StatusBadRequest, "Contact email is not valid.")
		return
	}
	if !hexColor.MatchString(s.PrimaryColor) {
		writeError(w, http.StatusBadRequest, "Primary colour must be hex like #4f46e5.")
		return
	}
	if !validLangs[s.DefaultLang] {
		s.DefaultLang = "de"
	}
	if err := h.Store.UpdateSite(r.Context(), &s); err != nil {
		serverError(w, r, err)
		return
	}
	site, err := h.Store.GetSite(r.Context())
	h.writeItemResult(w, r, site, err)
}

// PUT /api/admin/site/logo — multipart field "logo"; PNG, JPEG, WebP or SVG, ≤ 512 KB
func (h *handlers) adminUploadLogo(w http.ResponseWriter, r *http.Request) {
	data, ct, ok := readImageUpload(w, r, "logo", 512<<10, "image/png", "image/jpeg", "image/webp", "image/svg+xml")
	if !ok {
		return
	}
	if err := h.Store.SetSiteLogo(r.Context(), data, ct); err != nil {
		serverError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PUT /api/admin/site/portrait — multipart field "portrait"; JPEG, PNG or WebP, ≤ 1 MB
func (h *handlers) adminUploadPortrait(w http.ResponseWriter, r *http.Request) {
	data, ct, ok := readImageUpload(w, r, "portrait", 1<<20, "image/jpeg", "image/png", "image/webp")
	if !ok {
		return
	}
	if err := h.Store.SetPortrait(r.Context(), data, ct); err != nil {
		serverError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PUT /api/admin/site/home?lang=de|en — full HomeContent for that language;
// an empty object {} resets that language to its defaults.
func (h *handlers) adminUpdateHome(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if !validLangs[lang] {
		writeError(w, http.StatusBadRequest, "lang must be de or en")
		return
	}
	var in store.HomeContent
	if err := readJSON(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	reset := in.Headline == "" && len(in.Features) == 0 // {} → reset
	if !reset {
		in.Headline = strings.TrimSpace(in.Headline)
		if in.Headline == "" {
			writeError(w, http.StatusBadRequest, "Headline is required.")
			return
		}
		if len(in.Features) > 6 {
			writeError(w, http.StatusBadRequest, "At most 6 feature cards.")
			return
		}
		for i := range in.Features {
			f := &in.Features[i]
			f.Title, f.Text, f.Icon = strings.TrimSpace(f.Title), strings.TrimSpace(f.Text), strings.TrimSpace(f.Icon)
			if f.Title == "" {
				writeError(w, http.StatusBadRequest, "Each feature needs a title.")
				return
			}
			if !featureColors[f.Color] {
				f.Color = "teal"
			}
			if f.Icon == "" {
				f.Icon = "star_outline"
			}
		}
		in.FeaturedSlug = strings.ToLower(strings.TrimSpace(in.FeaturedSlug))
		if in.FeaturedSlug != "" && !validSlug(in.FeaturedSlug) {
			writeError(w, http.StatusBadRequest, "Featured slug is not a valid slug.")
			return
		}
	}

	if err := h.Store.UpdateHome(r.Context(), lang, in); err != nil {
		serverError(w, r, err)
		return
	}
	site, err := h.Store.GetSite(r.Context())
	h.writeItemResult(w, r, site, err)
}
