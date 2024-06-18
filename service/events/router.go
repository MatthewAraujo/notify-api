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
	router.HandleFunc("/events/{repoId}", h.getEventsByRepo).Methods(http.MethodGet)
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
	repoId, ok := vars["repoId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing repo name"))
		return
	}

	id, err := uuid.Parse(repoId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid repo id"))
		return
	}

	events, err := h.store.GetAllEventsForRepo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, events)
}
