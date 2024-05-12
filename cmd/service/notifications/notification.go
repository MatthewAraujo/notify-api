package notifications

import (
	"fmt"
	"net/http"

	"github.com/MatthewAraujo/notify/cmd/types"
	"github.com/MatthewAraujo/notify/cmd/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	// store types.NotificationStore
}

func NewHandler(
// store types.NotificationStore
) *Handler {
	return &Handler{
		// store: store
	}
}

func (h *Handler) Register(router *mux.Router) {
	router.HandleFunc("/notification", h.CreateNotification).Methods(http.MethodPost)

}

func (h *Handler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var payload types.Notifications
	if err := utils.ParseJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("validation error: %s", errors))
		return
	}

	err := CreateWebhook("MatthewAraujo", payload.RepoName, payload.Events)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, payload)
}
