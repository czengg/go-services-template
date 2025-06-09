package app

import (
	"log"
	"net/http"
	httphandlers "template/internal/adapters/inbound/http-handlers"
	"template/internal/logger"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type router struct {
	Upwardli httphandlers.UpwardliHandler
}

func newRouter(s services) router {
	return router{
		Upwardli: httphandlers.NewUpwardliHandler(s.upwardli),
	}
}

func (router *router) Serve(port string, logger logger.Logger) {
	r := mux.NewRouter()

	httphandlers.AcceptUpwardliEndpoints(r, router.Upwardli)

	corsOptions := []handlers.CORSOption{
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch, http.MethodOptions}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.MaxAge(3600),
	}

	logger.Info("Starting server...")
	log.Fatal(http.ListenAndServe(port, handlers.CORS(corsOptions...)(r)))
}
