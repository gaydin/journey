package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/dimfeld/httptreemux/v5"

	"github.com/gaydin/journey/filenames"
	"github.com/gaydin/journey/store"
	"github.com/gaydin/journey/structure/methods"
	"github.com/gaydin/journey/templates"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	number := httptreemux.ContextParams(r.Context())["number"]
	if number == "" {
		// Render index template (first page)
		err := templates.ShowIndexTemplate(r.Context(), store.DB, w, r, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	page, err := strconv.Atoi(number)
	if err != nil || page <= 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Render index template
	err = templates.ShowIndexTemplate(r.Context(), store.DB, w, r, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func authorHandler(w http.ResponseWriter, r *http.Request) {
	slug := httptreemux.ContextParams(r.Context())["slug"]
	function := httptreemux.ContextParams(r.Context())["function"]
	number := httptreemux.ContextParams(r.Context())["number"]
	if function == "" {
		// Render author template (first page)
		err := templates.ShowAuthorTemplate(r.Context(), store.DB, w, r, slug, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	} else if function == "rss" {
		// Render author rss feed
		err := templates.ShowAuthorRss(r.Context(), w, slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	page, err := strconv.Atoi(number)
	if err != nil || page <= 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Render author template
	err = templates.ShowAuthorTemplate(r.Context(), store.DB, w, r, slug, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func tagHandler(w http.ResponseWriter, r *http.Request) {
	slug := httptreemux.ContextParams(r.Context())["slug"]
	function := httptreemux.ContextParams(r.Context())["function"]
	number := httptreemux.ContextParams(r.Context())["number"]
	if function == "" {
		// Render tag template (first page)
		err := templates.ShowTagTemplate(r.Context(), store.DB, w, r, slug, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	} else if function == "rss" {
		// Render tag rss feed
		err := templates.ShowTagRss(r.Context(), w, slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	page, err := strconv.Atoi(number)
	if err != nil || page <= 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Render tag template
	err = templates.ShowTagTemplate(r.Context(), store.DB, w, r, slug, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	slug := httptreemux.ContextParams(r.Context())["slug"]
	if slug == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if slug == "rss" {
		// Render index rss feed
		err := templates.ShowIndexRss(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// Render post template
	err := templates.ShowPostTemplate(r.Context(), store.DB, w, r, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func newPostEditHandler(db store.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := httptreemux.ContextParams(r.Context())["slug"]
		if slug == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// Redirect to edit
		post, err := db.RetrievePostBySlug(r.Context(), slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		url := fmt.Sprintf("/admin/#/edit/%d", post.Id)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	// Read lock global blog
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	http.ServeFile(w, r, filepath.Join(filenames.ThemesFilepath, methods.Blog.ActiveTheme, "assets", httptreemux.ContextParams(r.Context())["filepath"]))
	return
}

func imagesHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(filenames.ImagesFilepath, httptreemux.ContextParams(r.Context())["filepath"]))
	return
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(filenames.PublicFilepath, httptreemux.ContextParams(r.Context())["filepath"]))
	return
}

func InitializeBlog(router *httptreemux.ContextMux, db store.Database) {
	// For index
	router.GET("/", indexHandler)
	router.GET("/:slug/edit", newPostEditHandler(db))
	router.GET("/:slug/", postHandler)
	router.GET("/page/:number/", indexHandler)
	// For author
	router.GET("/author/:slug/", authorHandler)
	router.GET("/author/:slug/:function/", authorHandler)
	router.GET("/author/:slug/:function/:number/", authorHandler)
	// For tag
	router.GET("/tag/:slug/", tagHandler)
	router.GET("/tag/:slug/:function/", tagHandler)
	router.GET("/tag/:slug/:function/:number/", tagHandler)
	// For serving asset files
	router.GET("/assets/*filepath", assetsHandler)
	router.GET("/images/*filepath", imagesHandler)
	router.GET("/content/images/*filepath", imagesHandler) // This is here to keep compatibility with Ghost
	router.GET("/public/*filepath", publicHandler)
}
