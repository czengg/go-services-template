package httphandlers

import (
	"net/http"
	"template/internal/core/upwardli"
)

type UpwardliHandler interface {
	CreateAllWebhooksHandler(w http.ResponseWriter, r *http.Request)
	CreateWebhookHandler(w http.ResponseWriter, r *http.Request)
	GetWebhooksHandler(w http.ResponseWriter, r *http.Request)
	DeleteWebhookHandler(w http.ResponseWriter, r *http.Request)
}

type upwardliHandler struct {
	service upwardli.Service
}

func NewUpwardliHandler(service upwardli.Service) UpwardliHandler {
	return &upwardliHandler{
		service: service,
	}
}

func (h *upwardliHandler) CreateAllWebhooksHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *upwardliHandler) CreateWebhookHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *upwardliHandler) GetWebhooksHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *upwardliHandler) DeleteWebhookHandler(w http.ResponseWriter, r *http.Request) {

}
