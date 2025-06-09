package router

import (
	"elevate/internal/handlers"
	"elevate/internal/logger"
	"elevate/internal/service"
	"log"
	"net/http"

	muxHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Router struct {
	Upwardli handlers.UpwardliHandler
}

func NewRouter(service service.Service) *Router {
	return &Router{
		Upwardli: handlers.NewUpwardliHandler(service),
	}
}

func (router *Router) Serve(cfg Config, logger logger.Logger) {
	r := mux.NewRouter()

	router.upwardliEndpoints(r)

	corsOptions := []muxHandlers.CORSOption{
		muxHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		muxHandlers.AllowedMethods([]string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch, http.MethodOptions}),
		muxHandlers.AllowedOrigins([]string{"*"}),
		muxHandlers.MaxAge(3600),
	}

	logger.Info("Starting server...")
	log.Fatal(http.ListenAndServe(cfg.Port, muxHandlers.CORS(corsOptions...)(r)))
}
