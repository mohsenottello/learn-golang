package httpapi

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// spaHandler serves the built React app. Paths that aren't real files fall
// back to index.html so client-side routing works on refresh.
func spaHandler(dist fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(dist))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if p == "" {
			p = "index.html"
		}
		if _, err := fs.Stat(dist, p); err != nil {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}
