package httphandlers

import (
	"io"
	"net/http"
	"strings"
	webhookprocessors "template/internal/adapters/inbound/webhook-processors"
	"template/internal/config"
	webhooks "template/internal/core/webhooks"
	"template/packages/common-go"
)

type UpwardliHandler interface {
	CreateAllWebhooksHandler(w http.ResponseWriter, r *http.Request)
	CreateWebhookHandler(w http.ResponseWriter, r *http.Request)
	GetWebhooksHandler(w http.ResponseWriter, r *http.Request)
	DeleteWebhookHandler(w http.ResponseWriter, r *http.Request)
	ProcessWebhookHandler(w http.ResponseWriter, r *http.Request)
}

type upwardliHandler struct {
	webhooksService webhooks.Service
	cfg             config.Config
	processor       webhooks.Processor
}

func NewUpwardliHandler(cfg config.Config, service webhooks.Service, processor webhooks.Processor) UpwardliHandler {
	return &upwardliHandler{
		webhooksService: service,
		cfg:             cfg,
		processor:       processor,
	}
}

func (h *upwardliHandler) CreateAllWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	err := h.webhooksService.CreateWebhooks(r.Context(), h.cfg.Upwardli().WebhookURL, []webhooks.SubscriptionTopic{
		webhookprocessors.SubscriptionTopicConsumerCreated,
		webhookprocessors.SubscriptionTopicConsumerUpdated,
		webhookprocessors.SubscriptionTopicConsumerClosed,
		webhookprocessors.SubscriptionTopicConsumerKYCStarted,
		webhookprocessors.SubscriptionTopicConsumerKYCPending,
		webhookprocessors.SubscriptionTopicConsumerKYCCompleted,
		webhookprocessors.SubscriptionTopicConsumerKYCNeedsReview,
		webhookprocessors.SubscriptionTopicConsumerKYCApproved,
		webhookprocessors.SubscriptionTopicConsumerKYCFailed,
		webhookprocessors.SubscriptionTopicPaymentCardCreated,
		webhookprocessors.SubscriptionTopicPaymentCardUpdated,
		webhookprocessors.SubscriptionTopicPaymentCardClosed,
		webhookprocessors.SubscriptionTopicPaymentCardTransactionSettlement,
		webhookprocessors.SubscriptionTopicACHSent,
		webhookprocessors.SubscriptionTopicACHReceived,
		webhookprocessors.SubscriptionTopicACHFailed,
		webhookprocessors.SubscriptionTopicPaymentTransferCreated,
		webhookprocessors.SubscriptionTopicPaymentTransferCompleted,
		webhookprocessors.SubscriptionTopicPaymentTransferFailed,
	})
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

func (h *upwardliHandler) ProcessWebhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	headers := make(map[string]string)
	for key, values := range r.Header {
		headers[key] = strings.Join(values, ",")
	}

	err = h.processor.Process(r.Context(), body, headers)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteJSON(w, http.StatusOK, "Webhook processed successfully")
}
