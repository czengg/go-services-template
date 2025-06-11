package httphandlers

import "github.com/go-chi/chi/v5"

func AcceptUpwardliEndpoints(r *chi.Mux, handler UpwardliHandler) {

	// TODO: authentication middleware belong here

	r.Route("/me/upwardli", func(r chi.Router) {
		r.Post("/webhooks", handler.CreateWebhookHandler)
		r.Post("/webhooks/all", handler.CreateAllWebhooksHandler)
		r.Get("/webhooks", handler.GetWebhooksHandler)
		r.Delete("/webhooks/{id}", handler.DeleteWebhookHandler)
	})

	r.Route("/admin/users/{userId}/upwardli", func(r chi.Router) {
		r.Post("/webhooks", handler.CreateWebhookHandler)
		r.Post("/webhooks/all", handler.CreateAllWebhooksHandler)
		r.Get("/webhooks", handler.GetWebhooksHandler)
		r.Delete("/webhooks/{id}", handler.DeleteWebhookHandler)
	})

}
