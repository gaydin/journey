package main

import (
	"net/http"
	"strings"

	"github.com/dimfeld/httptreemux/v5"
	"log/slog"

	"github.com/gaydin/journey/configuration"
	"github.com/gaydin/journey/database"
	"github.com/gaydin/journey/flags"
	"github.com/gaydin/journey/https"
	"github.com/gaydin/journey/logger"
	"github.com/gaydin/journey/server"
	"github.com/gaydin/journey/structure/methods"
	"github.com/gaydin/journey/templates"
)

func httpsRedirect(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	http.Redirect(w, r, configuration.Config.HttpsUrl+r.RequestURI, http.StatusMovedPermanently)
	return
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
	if err := database.Initialize(); err != nil {
		log.Error("Couldn't initialize database", logger.Error(err))
		return
	}

	// Global blog data
	if err := methods.GenerateBlog(); err != nil {
		log.Error("Couldn't generate blog data", logger.Error(err))
		return
	}

	// Templates
	if err := templates.Generate(); err != nil {
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
		httpsRouter := httptreemux.New()
		httpRouter := httptreemux.New()
		// Blog and pages as https
		server.InitializeBlog(httpsRouter)
		server.InitializePages(httpsRouter)
		// Add redirection to http router
		httpRouter.GET("/", httpsRedirect)
		httpRouter.GET("/*path", httpsRedirect)
		// Start https server
		log.Info("Starting https server on port " + httpsPort + "...")
		go func() {
			if err := https.StartServer(httpsPort, httpsRouter); err != nil {
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
		httpRouter := httptreemux.New()
		// Blog and pages as http
		server.InitializeBlog(httpRouter)
		server.InitializePages(httpRouter)
		// Start http server
		log.Info("Starting http server on port " + httpPort + "...")
		if err := http.ListenAndServe(httpPort, logger.Middleware(httpRouter, log)); err != nil {
			log.Error("Couldn't start the HTTP server:", logger.Error(err))
		}
	}
}
