package httpapi

// Serves the built React app. Because crawlers read the raw HTML, this handler
// injects per-page <title>/<meta>/<link> tags and sets a real HTTP status before
// handing the single-page app to the browser.

import (
	"encoding/json"
	"errors"
	"html"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/mohsenottello/golang/internal/store"
)

// App routes the SPA owns. Anything else is either a published resume slug or a 404.
var appRoutes = map[string]bool{
	"": true, "login": true, "register": true, "verify": true, "forgot-password": true,
	"reset-password": true, "account": true, "contact": true, "admin": true,
}

// Routes that must never appear in search results.
var noindexRoutes = map[string]bool{
	"login": true, "register": true, "verify": true, "forgot-password": true,
	"reset-password": true, "account": true, "admin": true,
}

type pageMeta struct {
	Title, Description, Image, OgType string
	Canonical                         string   // "" = none
	Alternates                        []string // hreflang links, already rendered
	NoIndex                           bool
	JSONLD                            any // nil = none
	Status                            int
}

func (h *handlers) spaHandler(dist fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(dist))
	index, _ := fs.ReadFile(dist, "index.html")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if p != "" && p != "index.html" {
			if _, err := fs.Stat(dist, p); err == nil {
				fileServer.ServeHTTP(w, r) // a real built asset
				return
			}
		}

		meta, err := h.metaFor(r, p)
		if err != nil {
			serverError(w, r, err)
			return
		}

		var head strings.Builder
		head.WriteString(`<meta property="og:title" content="` + meta.Title + `">`)
		head.WriteString(`<meta property="og:description" content="` + meta.Description + `">`)
		head.WriteString(`<meta property="og:type" content="` + meta.OgType + `">`)
		if meta.Image != "" {
			head.WriteString(`<meta property="og:image" content="` + meta.Image + `">`)
			head.WriteString(`<meta name="twitter:card" content="summary_large_image">`)
		}
		if meta.Canonical != "" {
			head.WriteString(`<link rel="canonical" href="` + meta.Canonical + `">`)
		}
		for _, a := range meta.Alternates {
			head.WriteString(a)
		}
		if meta.NoIndex {
			head.WriteString(`<meta name="robots" content="noindex">`)
		}
		if meta.JSONLD != nil {
			if b, err := json.Marshal(meta.JSONLD); err == nil {
				head.WriteString(`<script type="application/ld+json">` + string(b) + `</script>`)
			}
		}

		page := strings.NewReplacer(
			"__TITLE__", meta.Title,
			"__DESCRIPTION__", meta.Description,
			"</head>", head.String()+"</head>",
		).Replace(string(index))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(meta.Status)
		w.Write([]byte(page))
	})
}

// metaFor decides what a crawler should see for this path.
func (h *handlers) metaFor(r *http.Request, p string) (*pageMeta, error) {
	site, err := h.Store.GetSite(r.Context())
	if err != nil {
		return nil, err
	}
	base := strings.TrimSuffix(h.BaseURL, "/")
	segs := strings.Split(p, "/")
	first := segs[0]

	// --- app pages: homepage meta, private ones noindexed ---
	if appRoutes[first] {
		home := site.Home[site.DefaultLang]
		image := ""
		if site.HasPortrait {
			image = base + "/api/site/portrait"
		} else if site.HasLogo {
			image = base + "/api/site/logo"
		}
		m := &pageMeta{
			Title:       html.EscapeString(site.Name + " – " + site.TaglineEn),
			Description: html.EscapeString(home.Description),
			Image:       image,
			OgType:      "website",
			NoIndex:     noindexRoutes[first],
			Status:      http.StatusOK,
		}
		if first == "" {
			m.Canonical = base + "/"
			m.JSONLD = map[string]any{
				"@context": "https://schema.org", "@type": "WebSite",
				"name": site.Name, "url": base + "/",
				"publisher": map[string]any{
					"@type": "Organization", "name": site.Name, "url": base + "/",
					"logo": base + "/favicon-192.png",
				},
			}
		}
		return m, nil
	}

	// --- resume pages: /slug, /slug/de, /slug/en ---
	res, err := h.Store.ResumeBySlug(r.Context(), first)
	if errors.Is(err, store.ErrNotFound) || len(segs) > 2 {
		home := site.Home[site.DefaultLang]
		return &pageMeta{
			Title:       html.EscapeString(site.Name + " – 404"),
			Description: html.EscapeString(home.Description),
			OgType:      "website",
			NoIndex:     true,
			Status:      http.StatusNotFound,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	owner, err := h.Store.GetUser(r.Context(), res.UserID)
	if err != nil {
		return nil, err
	}

	lang := res.DefaultLang
	if len(segs) == 2 && validLangs[segs[1]] {
		lang = segs[1]
	}
	title := pick(lang, res.Title, res.TitleEn)
	about := pick(lang, res.About, res.AboutEn)
	url := base + "/" + res.Slug

	m := &pageMeta{
		Title:       html.EscapeString(owner.Name + " – " + title),
		Description: html.EscapeString(firstSentence(about, 160)),
		OgType:      "profile",
		Canonical:   url,
		Status:      http.StatusOK,
		Alternates: []string{
			`<link rel="alternate" hreflang="de" href="` + url + `/de">`,
			`<link rel="alternate" hreflang="en" href="` + url + `/en">`,
			`<link rel="alternate" hreflang="x-default" href="` + url + `">`,
		},
	}
	if res.HasPhoto {
		m.Image = base + "/api/resumes/" + res.Slug + "/photo"
	}
	person := map[string]any{
		"@context": "https://schema.org", "@type": "Person",
		"name": owner.Name, "jobTitle": title, "url": url,
	}
	if m.Image != "" {
		person["image"] = m.Image
	}
	if res.LinkedIn != "" {
		person["sameAs"] = []string{res.LinkedIn}
	}
	if res.Location != "" {
		person["address"] = map[string]any{"@type": "PostalAddress", "addressLocality": res.Location}
	}
	m.JSONLD = person
	return m, nil
}

// firstSentence trims a long text for a meta description.
func firstSentence(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	cut := s[:max]
	if i := strings.LastIndexAny(cut, ".!?"); i > max/2 {
		return cut[:i+1]
	}
	if i := strings.LastIndex(cut, " "); i > 0 {
		cut = cut[:i]
	}
	return cut + "…"
}
