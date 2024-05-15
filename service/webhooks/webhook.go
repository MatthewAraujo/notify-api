package webhooks

import (
	"fmt"
	"net/http"

	"github.com/MatthewAraujo/notify/service/mailer"
	"github.com/MatthewAraujo/notify/types"
	"github.com/MatthewAraujo/notify/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(mux *mux.Router) {
	mux.HandleFunc("/webhooks", h.webhooksHandler).Methods(http.MethodPost)
	mux.HandleFunc("/webhooks/installation", h.installationHandler).Methods(http.MethodPost)
}

func (h *Handler) installationHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.InstallationWebhooks
	if err := utils.ParseJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("validation error: %s", errors))
		return
	}

	//store installation ID for the user
	//store repositories for the user
}

func (h *Handler) webhooksHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.GithubWebhooks
	if err := utils.ParseJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("validation error: %s", errors))
		return
	}

	bodyEmail := types.SendEmail{
		RepoName: payload.Repository.FullName,
		Sender:   payload.Repository.Owner.Name,
		Commit:   payload.Commits[0].Message,
		Email:    payload.Repository.Owner.Email,
	}

	go mailer.SendMail(bodyEmail)
}
