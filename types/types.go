package types

import (
	"time"

	"github.com/google/uuid"
)

type InstallationWebhooks struct {
	Installation struct {
		Id      int `json:"id"`
		Account struct {
			Login string `json:"login"`
		}
	} `json:"installation"`
	Repositories []struct {
		Name string `json:"name"`
	}
}

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

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Repository struct {
	ID        uuid.UUID `json:"id"`
	RepoName  string    `json:"repo_name"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Installation struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"` // pointer to time.Time to allow nil values
	CreatedAt time.Time  `json:"created_at"`
}

type EventType struct {
	ID        uuid.UUID `json:"id"`
	EventName string    `json:"event_name"`
}

type Event struct {
	ID        uuid.UUID `json:"id"`
	RepoID    uuid.UUID `json:"repo_id"`
	EventType uuid.UUID `json:"event_type"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationSubscription struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	RepoID    uuid.UUID `json:"repo_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
