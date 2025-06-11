package app

import (
	"log"
	"net/http"
	httphandlers "template/internal/adapters/inbound/http-handlers"
	"template/internal/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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
	r := chi.NewRouter()

	httphandlers.AcceptUpwardliEndpoints(r, router.Upwardli)

	r.Use(cors.Handler(cors.Options{
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		AllowedMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowedOrigins: []string{"*"},
		MaxAge:         3600,
	}))

	logger.Info("Starting server...")
	log.Fatal(http.ListenAndServe(port, r))
}
