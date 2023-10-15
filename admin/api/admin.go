package admin

import (
	"net/http"

	"github.com/gaydin/journey/admin/api/cors"
	"github.com/gaydin/journey/admin/api/security"
	"github.com/gaydin/journey/admin/generated"
	"github.com/gaydin/journey/store"
)

type Handler struct {
	db  store.Database
	url string
}

func New(url string, db store.Database) http.Handler {
	apiHandlers := &Handler{
		db:  db,
		url: url,
	}

	securityHandler := security.New(db)

	server, err := generated.NewServer(apiHandlers, securityHandler)
	if err != nil {
		panic(err)
	}

	return cors.NewMiddleware(server)
}
