package server

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dimfeld/httptreemux/v5"

	"github.com/gaydin/journey/filenames"
	"github.com/gaydin/journey/helpers"
)

func pagesHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(filenames.PagesFilepath, httptreemux.ContextParams(r.Context())["filepath"])
	// If the path points to a directory, add a trailing slash to the path (needed if the page loads relative assets).
	if helpers.IsDirectory(path) && !strings.HasSuffix(r.RequestURI, "/") {
		http.Redirect(w, r, r.RequestURI+"/", 301)
		return
	}
	http.ServeFile(w, r, path)
	return
}

func InitializePages(router *httptreemux.ContextMux) {
	// For serving standalone projects or pages saved in content/pages
	router.GET("/pages/*filepath", pagesHandler)
}
