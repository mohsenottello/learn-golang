package httpapi

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// readImageUpload parses a single-file multipart upload and checks size and type.
// On failure it has already written the error response and returns ok=false.
func readImageUpload(w http.ResponseWriter, r *http.Request, field string, maxBytes int64, allowed ...string) ([]byte, string, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes+4096)
	if err := r.ParseMultipartForm(maxBytes); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Image must be smaller than %d KB.", maxBytes>>10))
		return nil, "", false
	}
	f, _, err := r.FormFile(field)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing '"+field+"' field")
		return nil, "", false
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not read file")
		return nil, "", false
	}
	ct := http.DetectContentType(data)
	if strings.HasPrefix(ct, "text/") && strings.Contains(strings.ToLower(string(data[:min(len(data), 512)])), "<svg") {
		ct = "image/svg+xml"
	}
	for _, a := range allowed {
		if ct == a {
			return data, ct, true
		}
	}
	writeError(w, http.StatusBadRequest, "Unsupported image type.")
	return nil, "", false
}
