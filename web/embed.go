// Package web embeds the built React app so the Go binary is self-contained.
package web

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var distFS embed.FS

// Dist is the dist/ folder as a filesystem rooted at its contents.
var Dist, _ = fs.Sub(distFS, "dist")
