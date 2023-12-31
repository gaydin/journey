package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/dimfeld/httptreemux/v5"
	"log/slog"

	"github.com/gaydin/journey/admin"
	"github.com/gaydin/journey/configuration"
	"github.com/gaydin/journey/flags"
	"github.com/gaydin/journey/https"
	"github.com/gaydin/journey/logger"
	"github.com/gaydin/journey/server"
	"github.com/gaydin/journey/store"
	"github.com/gaydin/journey/structure/methods"
	"github.com/gaydin/journey/templates"
)

func newHttpRedirect(httpsURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, httpsURL+r.RequestURI, http.StatusMovedPermanently)
	}
}

func main() {
	// Configuration is read from config.json by loading the configuration package
	config := configuration.NewConfiguration()

	log, closeLogFunc := logger.New(config)
	defer func() {
		if err := closeLogFunc(); err != nil {
			slog.Default().Error("close log file error", logger.Error(err))
		}
	}()

	// Database
	db, err := store.New(config)
	if err != nil {
		log.Error("Couldn't initialize database", logger.Error(err))
		return
	}

	// Global blog data
	if err := methods.GenerateBlog(context.Background(), config.Url, db); err != nil {
		log.Error("Couldn't generate blog data", logger.Error(err))
		return
	}

	// Templates
	if err := templates.Generate(context.Background(), db); err != nil {
		log.Error("Couldn't compile templates", logger.Error(err))
		return
	}

	// HTTP(S) Server
	httpPort := configuration.Config.HttpHostAndPort
	httpsPort := configuration.Config.HttpsHostAndPort
	// Check if HTTP/HTTPS flags were provided
	if flags.HttpPort != "" {
		components := strings.SplitAfterN(httpPort, ":", 2)
		httpPort = components[0] + flags.HttpPort
	}
	if flags.HttpsPort != "" {
		components := strings.SplitAfterN(httpsPort, ":", 2)
		httpsPort = components[0] + flags.HttpsPort
	}

	// Determine the kind of https support (as set in the config.json)
	if configuration.Config.HttpsUsage {
		httpsRouter := httptreemux.NewContextMux()
		httpRouter := httptreemux.NewContextMux()
		// Blog and pages as https
		server.InitializeBlog(httpsRouter, db)
		server.InitializePages(httpsRouter)
		admin.InitializeAdmin(config, db, httpsRouter)
		// Add redirection to http router
		httpRouter.GET("/", newHttpRedirect(config.HttpsUrl))
		httpRouter.GET("/*path", newHttpRedirect(config.HttpsUrl))
		// Start https server
		log.Info("Starting https server on port " + httpsPort + "...")
		go func() {
			if err := https.StartServer(httpsPort, logger.Middleware(httpsRouter, log)); err != nil {
				log.Error("Couldn't start the HTTPS server", logger.Error(err))
				return
			}
		}()
		// Start http server
		log.Info("Starting http server on port " + httpPort + "...")
		if err := http.ListenAndServe(httpPort, logger.Middleware(httpRouter, log)); err != nil {
			log.Error("Couldn't start the HTTP server:", logger.Error(err))
		}
	} else {
		httpRouter := httptreemux.NewContextMux()
		// Blog and pages as http
		server.InitializeBlog(httpRouter, db)
		server.InitializePages(httpRouter)
		admin.InitializeAdmin(config, db, httpRouter)

		// Start http server
		log.Info("Starting http server on port " + httpPort + "...")
		if err := http.ListenAndServe(httpPort, logger.Middleware(httpRouter, log)); err != nil {
			log.Error("Couldn't start the HTTP server:", logger.Error(err))
		}
	}
}
