package httphandlers

import (
	"net/http"
	"template/internal/config"
	webhooks "template/internal/core/webhooks"
	"template/packages/common-go"
)

type UpwardliHandler interface {
	CreateAllWebhooksHandler(w http.ResponseWriter, r *http.Request)
	CreateWebhookHandler(w http.ResponseWriter, r *http.Request)
	GetWebhooksHandler(w http.ResponseWriter, r *http.Request)
	DeleteWebhookHandler(w http.ResponseWriter, r *http.Request)
}

type upwardliHandler struct {
	webhooksService webhooks.Service
	cfg             config.Config
}

func NewUpwardliHandler(cfg config.Config, service webhooks.Service) UpwardliHandler {
	return &upwardliHandler{
		webhooksService: service,
		cfg:             cfg,
	}
}

func (h *upwardliHandler) CreateAllWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	err := h.webhooksService.CreateAllWebhooks(r.Context(), h.cfg.Upwardli().WebhookURL)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteJSON(w, http.StatusOK, "Webhooks created successfully")
}

func (h *upwardliHandler) CreateWebhookHandler(w http.ResponseWriter, r *http.Request) {
	err := h.webhooksService.CreateWebhook(r.Context(),
		r.URL.Query().Get("endpoint"),
		webhooks.SubscriptionTopic(r.URL.Query().Get("webhookName")),
	)

	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteJSON(w, http.StatusOK, "Webhook created successfully")
}

func (h *upwardliHandler) GetWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	webhooks, err := h.webhooksService.GetWebhooks(r.Context())
	if err != nil {
		common.WriteError(w, err)
		return
	}

	response := make([]WebhookResponse, len(webhooks))
	for i, webhook := range webhooks {
		response[i] = WebhookToResponse(webhook)
	}

	common.WriteJSON(w, http.StatusOK, response)
}

func (h *upwardliHandler) DeleteWebhookHandler(w http.ResponseWriter, r *http.Request) {
	err := h.webhooksService.DeleteWebhook(r.Context(), r.URL.Query().Get("id"))
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteJSON(w, http.StatusOK, "Webhook deleted successfully")
}
