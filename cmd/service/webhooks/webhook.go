package webhooks

import (
	"net/http"

	"github.com/MatthewAraujo/notify/cmd/service/mailer"
	"github.com/MatthewAraujo/notify/cmd/types"
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
	var payload types.GithubWebhooks
	if err := utils.ParseJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
