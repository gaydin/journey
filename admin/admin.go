package admin

import (
	"github.com/dimfeld/httptreemux/v5"

	admin "github.com/gaydin/journey/admin/api"
	"github.com/gaydin/journey/configuration"
	"github.com/gaydin/journey/store"
)

func InitializeAdmin(
	config *configuration.Configuration,
	db store.Database,
	router *httptreemux.ContextMux) {

	handler := admin.New(config.Url, db)
	router.UsingContext().GET("/admin/v1/*path", handler.ServeHTTP)
	router.UsingContext().POST("/admin/v1/*path", handler.ServeHTTP)
	router.UsingContext().PATCH("/admin/v1/*path", handler.ServeHTTP)
	router.UsingContext().DELETE("/admin/v1/*path", handler.ServeHTTP)
}
