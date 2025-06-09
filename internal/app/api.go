package app

import (
	"log"
	"net/http"
	inboundHTTP "template/internal/adapters/inbound/http"

	"template/internal/logger"

	muxHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type router struct {
	Upwardli inboundHTTP.UpwardliHandler
}

func newRouter(s services) *router {
	return &router{
		Upwardli: inboundHTTP.NewUpwardliHandler(s.upwardli),
	}
}

func (router *router) Serve(port string, logger logger.Logger) {
	r := mux.NewRouter()

	inboundHTTP.AcceptUpwardliEndpoints(r, router.Upwardli)

	corsOptions := []muxHandlers.CORSOption{
		muxHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		muxHandlers.AllowedMethods([]string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch, http.MethodOptions}),
		muxHandlers.AllowedOrigins([]string{"*"}),
		muxHandlers.MaxAge(3600),
	}

	logger.Info("Starting server...")
	log.Fatal(http.ListenAndServe(port, muxHandlers.CORS(corsOptions...)(r)))
}
