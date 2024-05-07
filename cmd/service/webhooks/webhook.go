package webhooks

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(mux *mux.Router) {
	mux.HandleFunc("/webhooks", h.webhooksHandler)
}

func (h *Handler) webhooksHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
