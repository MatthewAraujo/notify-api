package webhooks

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/MatthewAraujo/notify/service/mailer"
	"github.com/MatthewAraujo/notify/service/notifications"
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
	var payload types.GithubInstallation
	if err := utils.ParseJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if payload.Action == "created" {
		if err := utils.Validate.Struct(payload); err != nil {
			errors := err.(validator.ValidationErrors)
			log.Printf("validation error: %s", errors)
			return
		}

		user, err := h.store.GetUserIdByUsername(payload.Installation.Account.Login)
		if err != nil && err.Error() == "user not found" {
			log.Printf("User not found for %s", payload.Installation.Account.Login)
			return
		}

		installationId := payload.Installation.Id

		//Check if the installation already exists
		exists, err := h.store.CheckIfInstallationExists(user.ID)
		if err != nil {
			return
		}

		if exists {
			log.Printf("Installation already exists for %s", payload.Installation.Account.Login)
			return
		}

		if err := h.store.CreateInstallation(user.ID, installationId); err != nil {
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

			if err := h.store.CreateRepository(user.ID, repo.Name); err != nil {
				log.Printf("Error creating repository for %s", repo.Name)
				return
			}
		}

		utils.WriteJSON(w, http.StatusCreated, fmt.Sprintf("Installation created for %s", payload.Installation.Account.Login))
	}
	if payload.Action == "added" {

		if err := utils.Validate.Struct(payload); err != nil {
			errors := err.(validator.ValidationErrors)
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("validation error: %s", errors))
			return
		}

		userId, err := h.store.GetUserIdByInstallationId(payload.Installation.Id)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		for _, repo := range payload.RepositoriesAdded {
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
	}

	if payload.Action == "deleted" {

		if err := utils.Validate.Struct(payload); err != nil {
			return
		}

		userId, err := h.store.GetUserIdByInstallationId(payload.Installation.Id)
		if err != nil {
			return
		}

		log.Printf("Removing subscription for %s", userId)

		// remove on github
		err = notifications.DeleteAllWebhooks(userId, h.store)
		if err != nil {
			return
		}

		err = h.store.RevokeUser(userId)
		if err != nil {
			return
		}

	} else {
		return
	}
}

func (h *Handler) webhooksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received webhook")
	var payload types.GithubWebhooks
	if err := utils.ParseJSON(r, &payload); err != nil {
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		return
	}

	hook, err := h.store.CheckIfHookIdExistsInNotificationSubscription(payload.HookId)
	if err != nil {
		return
	}

	log.Printf("Hook exists: %t", hook)
	log.Printf("Received webhook for %s", payload.Repository.Owner.Name)

	if !hook {
		err := h.store.AddHookIdInNotificationSubscription(payload.Repository.Name, payload.HookId)
		if err != nil {
			if err.Error() == "repo not found" {
				return
			}
			return
		}
	}

	if strings.Contains(payload.Ref, "refs/heads/") {
		user, err := h.store.GetUserIdByUsername(payload.Repository.Owner.Name)
		if err != nil {
			return
		}

		bodyEmail := types.SendEmail{
			RepoName: payload.Repository.FullName,
			Sender:   payload.Repository.Owner.Name,
			Commit:   payload.Commits[0].Message,
			Email:    user.Email,
		}

		go mailer.SendMail(bodyEmail)
	}

	bodyEmail := types.WelcomeEmail{
		Email:      payload.Repository.Owner.Email,
		Owner:      payload.Repository.Owner.Name,
		Repository: payload.Repository.FullName,
	}

	go mailer.SendWelcomeEmail(bodyEmail)

}
