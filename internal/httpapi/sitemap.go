package httpapi

import (
	"net/http"
	"strings"
)

// GET /sitemap.xml
func (h *handlers) sitemap(w http.ResponseWriter, r *http.Request) {
	list, err := h.Store.ListResumesForAdmin(r.Context())
	if err != nil {
		serverError(w, r, err)
		return
	}
	base := strings.TrimSuffix(h.BaseURL, "/")
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	b.WriteString("<url><loc>")
	b.WriteString(base)
	b.WriteString("/</loc></url><url><loc>")
	b.WriteString(base)
	b.WriteString("/contact</loc></url>")
	for _, x := range list {
		if x.Published {
			b.WriteString("<url><loc>")
			b.WriteString(base)
			b.WriteString("/")
			b.WriteString(x.Slug)
			b.WriteString("</loc><lastmod>")
			b.WriteString(x.UpdatedAt.Format("2006-01-02"))
			b.WriteString("</lastmod></url>")
		}
	}
	b.WriteString("</urlset>")
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(b.String()))
}
