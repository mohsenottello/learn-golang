package httpapi

import (
	"html"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// spaHandler serves the built React app. Paths that aren't real files fall
// back to index.html so client-side routing works on refresh.
func (h *handlers) spaHandler(dist fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(dist))
	index, _ := fs.ReadFile(dist, "index.html")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if p != "" && p != "index.html" {
			if _, err := fs.Stat(dist, p); err == nil {
				fileServer.ServeHTTP(w, r) // real asset
				return
			}
		}
		site, err := h.Store.GetSite(r.Context())
		if err != nil {
			serverError(w, r, err)
			return
		}
		base := strings.TrimSuffix(h.BaseURL, "/")
		home := site.Home[site.DefaultLang]
		image := base + "/api/site/portrait"
		if !site.HasPortrait && site.HasLogo {
			image = base + "/api/site/logo"
		}

		title := html.EscapeString(site.Name + " – " + home.Headline)
		desc := html.EscapeString(home.Description)
		og := `<meta property="og:title" content="` + title + `">` +
			`<meta property="og:description" content="` + desc + `">` +
			`<meta property="og:image" content="` + image + `">` +
			`<meta property="og:type" content="website">` +
			`<meta name="twitter:card" content="summary_large_image">`
		page := strings.NewReplacer(
			"__TITLE__", title,
			"__DESCRIPTION__", desc,
			"</head>", og+"</head>",
		).Replace(string(index))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Write([]byte(page))
	})
}
