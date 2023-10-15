package admin

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"

	admin "github.com/gaydin/journey/admin/api"
	"github.com/gaydin/journey/admin/vue"
	"github.com/gaydin/journey/configuration"
	"github.com/gaydin/journey/store"
)

func InitializeAdmin(
	config *configuration.Configuration,
	db store.Database,
	router *httptreemux.ContextMux) {

	v := vue.Handler()

	router.GET("/admin/*filepath", v)
	router.GET("/admin/", newAdminHandler(db, v))

	handler := admin.New(config.Url, db)
	router.UsingContext().GET("/admin/v1/*path", handler.ServeHTTP)
	router.UsingContext().POST("/admin/v1/*path", handler.ServeHTTP)
	router.UsingContext().PATCH("/admin/v1/*path", handler.ServeHTTP)
	router.UsingContext().DELETE("/admin/v1/*path", handler.ServeHTTP)
}

// Function to route the /admin/ url accordingly. (Is user logged in? Is at least one user registered?)
func newAdminHandler(db store.Database, fh http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db.RetrieveUsersCount(r.Context()) == 0 {
			http.Redirect(w, r, "/admin/register/", http.StatusFound)
			return
		} else {
			fh.ServeHTTP(w, r)
		}
	}
}
