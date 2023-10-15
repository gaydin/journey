package vue

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var FS embed.FS

func Handler() http.HandlerFunc {
	// Load the files subdirectory
	fsys, err := fs.Sub(FS, "dist")
	if err != nil {
		panic(err)
	}

	// Create an http.FileServer to serve the
	// contents of the files subdiretory.
	handler := http.StripPrefix("/admin/", http.FileServer(http.FS(fsys)))

	// Create an http.HandlerFunc that wraps the
	// http.FileServer to always load the index.html
	// file if a directory path is being requested.
	return func(w http.ResponseWriter, r *http.Request) {

		handler.ServeHTTP(w, r)
	}
}
