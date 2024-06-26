package health

import (
	"net/http"

	"github.com/MatthewAraujo/notify/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(mux *mux.Router) {
	mux.HandleFunc("/health", h.healthHandler)
}

func (h *Handler) healthHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, "api is on")
}
