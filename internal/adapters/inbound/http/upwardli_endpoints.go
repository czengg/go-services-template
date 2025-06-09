package http

import (
	"github.com/gorilla/mux"
)

func AcceptUpwardliEndpoints(r *mux.Router, handler UpwardliHandler) {
	upwardliUserRouter := r.PathPrefix("/me/upwardli").Subrouter()
	upwardliUserRouter.HandleFunc("/webhooks", handler.CreateWebhookHandler).Methods("POST")
	upwardliUserRouter.HandleFunc("/webhooks/all", handler.CreateAllWebhooksHandler).Methods("POST")
	upwardliUserRouter.HandleFunc("/webhooks", handler.GetWebhooksHandler).Methods("GET")
	upwardliUserRouter.HandleFunc("/webhooks/{id}", handler.DeleteWebhookHandler).Methods("DELETE")

	upwardliAdminRouter := r.PathPrefix("/admin/users/{userId}/upwardli").Subrouter()
	upwardliAdminRouter.HandleFunc("/webhooks", handler.CreateWebhookHandler).Methods("POST")
	upwardliAdminRouter.HandleFunc("/webhooks/all", handler.CreateAllWebhooksHandler).Methods("POST")
	upwardliAdminRouter.HandleFunc("/webhooks", handler.GetWebhooksHandler).Methods("GET")
	upwardliAdminRouter.HandleFunc("/webhooks/{id}", handler.DeleteWebhookHandler).Methods("DELETE")
}
