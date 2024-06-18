package repository

import (
	"fmt"
	"net/http"

	"github.com/MatthewAraujo/notify/types"
	"github.com/MatthewAraujo/notify/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.RepositoryStore
}

func NewHandler(store types.RepositoryStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) Register(router *mux.Router) {
	router.HandleFunc("/repository/{username}", h.GetAllRepositoryForUser).Methods(http.MethodGet)
}

func (h *Handler) GetAllRepositoryForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing id"))
		return
	}

	// check repos,
	// check if repo is a NotificationSubscription
	// get all events for that repoo

	repos, err := h.store.GetAllReposForUser(username)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if len(repos) == 0 {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("repository not found"))
		return
	}

	reposWithEvents := make([]types.ReposWithEvents, 0)

	for _, repo := range repos {
		isSubscribed, err := h.store.IsRepoSubscribed(username, repo.ID)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if isSubscribed {
			repos, err := h.store.GetAllRepositoryForUser(username)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			reposWithEvents = append(reposWithEvents, repos...)
		} else {
			reposWithEvents = append(reposWithEvents, types.ReposWithEvents{
				RepoId:   repo.ID,
				RepoName: repo.RepoName,
			})
		}
	}

	if len(reposWithEvents) == 0 {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("no repositories found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, reposWithEvents)

}
