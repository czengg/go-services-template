package router

import "github.com/gorilla/mux"

func (router *Router) upwardliEndpoints(r *mux.Router) {

	upwardliUserRouter := r.PathPrefix("/me/upwardli").Subrouter()
	upwardliUserRouter.HandleFunc("/webhooks", router.Upwardli.CreateWebhookHandler).Methods("POST")
	upwardliUserRouter.HandleFunc("/webhooks/all", router.Upwardli.CreateAllWebhooksHandler).Methods("POST")
	upwardliUserRouter.HandleFunc("/webhooks", router.Upwardli.GetWebhooksHandler).Methods("GET")
	upwardliUserRouter.HandleFunc("/webhooks/{id}", router.Upwardli.DeleteWebhookHandler).Methods("DELETE")

	upwardliAdminRouter := r.PathPrefix("/admin/users/{userId}/upwardli").Subrouter()
	upwardliAdminRouter.HandleFunc("/webhooks", router.Upwardli.CreateWebhookHandler).Methods("POST")
	upwardliAdminRouter.HandleFunc("/webhooks/all", router.Upwardli.CreateAllWebhooksHandler).Methods("POST")
	upwardliAdminRouter.HandleFunc("/webhooks", router.Upwardli.GetWebhooksHandler).Methods("GET")
	upwardliAdminRouter.HandleFunc("/webhooks/{id}", router.Upwardli.DeleteWebhookHandler).Methods("DELETE")
}
