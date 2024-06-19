package events

import (
	"fmt"
	"net/http"

	"github.com/MatthewAraujo/notify/types"
	"github.com/MatthewAraujo/notify/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.EventStore
}

func NewHandler(store types.EventStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) Register(router *mux.Router) {
	router.HandleFunc("/events", h.getAllEvents).Methods(http.MethodGet)
	router.HandleFunc("/events/{reponame}", h.getEventsByRepo).Methods(http.MethodGet)
}

func (h *Handler) getAllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.store.GetAllEvents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, events)
}

func (h *Handler) getEventsByRepo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName, ok := vars["reponame"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing repo name"))
		return
	}

	userId := h.store.GetUserIDFromRepoName(repoName)
	if userId == uuid.Nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("repo not found"))
		return
	}

	notificationSubscriptionId, err := h.store.GetNotificationSubscriptionId(userId, repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	events, err := h.store.GetAllEventsForRepo(repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(events) == 0 {
		payload := struct {
			UserID   uuid.UUID `json:"user_id"`
			Reponame string    `json:"reponame"`
		}{}
		payload.UserID = userId
		payload.Reponame = repoName

		utils.WriteJSON(w, http.StatusOK, payload)
	} else {
		userWithEvents := struct {
			UserID                     uuid.UUID         `json:"user_id"`
			Events                     []types.EventType `json:"events"`
			Reponame                   string            `json:"reponame"`
			NotificationSubscriptionId uuid.UUID         `json:"notification_subscription_id"`
		}{}
		userWithEvents.UserID = userId
		userWithEvents.Events = events
		userWithEvents.Reponame = repoName
		userWithEvents.NotificationSubscriptionId = notificationSubscriptionId

		utils.WriteJSON(w, http.StatusOK, userWithEvents)
	}
}
