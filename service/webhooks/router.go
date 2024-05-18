package webhooks

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MatthewAraujo/notify/service/mailer"
	"github.com/MatthewAraujo/notify/types"
	"github.com/MatthewAraujo/notify/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.InstallationStore
}

func NewHandler(store types.InstallationStore) *Handler {
	return &Handler{
		store: store,
	}
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

	userId, err := h.store.GetUserIdByUsername(payload.Installation.Account.Login)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	installationId := payload.Installation.Id

	//Check if the installation already exists
	exists, err := h.store.CheckIfInstallationExists(userId)
	if err != nil {
		return
	}

	if exists {
		log.Printf("Installation already exists for %s", payload.Installation.Account.Login)
		return
	}

	if err := h.store.CreateInstallation(userId, installationId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	//store repositories for the user
	for _, repo := range payload.Repositories {
		//Check if the repo already exists
		exists, err := h.store.CheckIfRepoExists(repo.Name)
		if err != nil {
			return
		}

		if exists {
			log.Printf("Repository already exists for %s", repo.Name)
			return
		}

		if err := h.store.CreateRepository(userId, repo.Name); err != nil {
			log.Printf("Error creating repository for %s", repo.Name)
			return
		}
	}

	utils.WriteJSON(w, http.StatusCreated, fmt.Sprintf("Installation created for %s", payload.Installation.Account.Login))
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
