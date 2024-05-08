package webhooks

import (
	"fmt"
	"net/http"

	"github.com/MatthewAraujo/notify/cmd/utils"
	"github.com/gorilla/mux"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(mux *mux.Router) {
	mux.HandleFunc("/webhooks", h.webhooksHandler).Methods("POST")
}

func (h *Handler) webhooksHandler(w http.ResponseWriter, r *http.Request) {
	var payload any
	if err := utils.ParseJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Print(payload)

}
