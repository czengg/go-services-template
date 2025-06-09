package handlers

import (
	"elevate/internal/service"
	"net/http"
)

type UpwardliHandler interface {
	CreateAllWebhooksHandler(w http.ResponseWriter, r *http.Request)
	CreateWebhookHandler(w http.ResponseWriter, r *http.Request)
	GetWebhooksHandler(w http.ResponseWriter, r *http.Request)
	DeleteWebhookHandler(w http.ResponseWriter, r *http.Request)
}

type upwardliHandler struct {
	service service.Service
}

func NewUpwardliHandler(service service.Service) UpwardliHandler {
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
