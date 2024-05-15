package types

import "github.com/google/uuid"

type GithubWebhooks struct {
	Repository struct {
		FullName string `json:"full_name"`
		Owner    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"owner"`
	} `json:"repository"`
	Commits []struct {
		Message string `json:"message"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			// Outros campos do autor do commit, se necessário
		} `json:"author"`
	} `json:"commits"`
}

type SendEmail struct {
	RepoName string
	Sender   string
	Commit   string
	Email    string
}

type NotificationStore interface {
	GetUserByID(id uuid.UUID) (*Notifications, error)
	GetRepositoryByUserID(id uuid.UUID, reponame string) (*Notifications, error)
	CreateNotification(notif *Notifications) error
}
type Notifications struct {
	UserId   uuid.UUID `json:"user_id"`
	RepoName string    `json:"repository_name"`
	Events   []string  `json:"events"`
}
