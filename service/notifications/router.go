package notifications

import (
	"fmt"
	"net/http"

	"github.com/MatthewAraujo/notify/types"
	"github.com/MatthewAraujo/notify/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.NotificationStore
}

func NewHandler(store types.NotificationStore) *Handler {
	return &Handler{
		store: store,
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

	user, err := h.store.GetUserByID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user not found"))
		return
	}

	installationId, err := h.store.GetInstallationIDByUser(user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	for _, repo := range payload.Repos {
		// check if the repo already exists
		exists, err := h.store.CheckIfRepoExists(repo.RepoName)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if !exists {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("repo not found"))
			return
		}
		repoId, err := h.store.GetRepoIDByName(repo.RepoName)
		if err != nil {
			if err.Error() == "repo not found" {
				utils.WriteError(w, http.StatusNotFound, err)
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		// check if the notification already exists
		exists, err = h.store.CheckIfNotificationExists(user.ID, repoId)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if exists {
			utils.WriteError(w, http.StatusConflict, fmt.Errorf("notification already exists"))
			return
		}

		err = CreateWebhook(installationId, user.Username, user.ID, repo.RepoName, repo.Events)

		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		for _, event := range repo.Events {
			eventId, err := h.store.GetEventTypeByName(event)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			notif := &types.NotificationSubscription{
				UserID: user.ID,
				RepoID: repoId,
			}

			if err := h.store.CreateNotification(notif); err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			event := &types.Event{
				RepoID:    repoId,
				EventType: eventId,
			}

			if err := h.store.CreateEvent(event); err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}
		}
	}

	utils.WriteJSON(w, http.StatusCreated, fmt.Sprintf("Notification created for %s", user.Username))
}
