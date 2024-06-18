package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MatthewAraujo/notify/types"
	"github.com/MatthewAraujo/notify/utils"
	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) Register(router *mux.Router) {
	router.HandleFunc("/delete", h.deleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/auth/{provider}/callback", h.getAuthCallbackFunction).Methods(http.MethodGet)
	router.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler).Methods(http.MethodGet)
	router.HandleFunc("/logout/{provider}", h.logout).Methods(http.MethodGet)

	// Get user infos
	router.HandleFunc("/user/{username}", h.getUser).Methods(http.MethodGet)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing username"))
		return
	}

	user, err := h.store.GetUserByUsername(username)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)

}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	var payload types.User
	if err := utils.ParseJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user not found"))
		return
	}

	if u.Email != payload.Email {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user not found"))
		return
	}

	err = h.store.DeleteUser(u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "user deleted"})
}

func (s *Handler) getAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	type contextKey string

	const providerKey contextKey = "provider"

	provider, ok := vars["provider"]
	if !ok {
		http.Error(w, "missing provider", http.StatusBadRequest)
		return
	}

	r = r.WithContext(context.WithValue(context.Background(), providerKey, provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.store.CreateUser(&types.User{
		Username:  user.NickName,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expiration := time.Now().Add(30 * 24 * time.Hour)
	cookie := http.Cookie{Name: "username", Value: user.NickName, Path: "/", Expires: expiration, HttpOnly: true}

	http.SetCookie(w, &cookie)
	urlRedirect := "http://localhost:3000/installation"
	http.Redirect(w, r, urlRedirect, http.StatusFound)

}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	http.Redirect(w, r, "http://localhost:3000", http.StatusFound)

}
