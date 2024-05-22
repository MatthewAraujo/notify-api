package types

import (
	"time"

	"github.com/google/uuid"
)

type GithubInstallation struct {
	Action       string `json:"action"`
	Installation struct {
		Id      int `json:"id"`
		Account struct {
			Login string `json:"login"`
		}
	} `json:"installation"`
	Repositories        []Repos `json:"repositories"`
	RepositoriesAdded   []Repos `json:"repositories_added"`
	RepositoriesRemoved []Repos `json:"repositories_removed"`
}
type Repos struct {
	Name string `json:"name"`
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
	// Notification
	CheckIfNotificationExists(id uuid.UUID) (bool, error)
	CheckIfNotificationExistsForUserId(userId uuid.UUID, repoId uuid.UUID) (bool, error)
	CreateNotification(notif *NotificationSubscription) error
	GetOwnerOfNotification(id uuid.UUID) (uuid.UUID, error)
	DeleteNotification(id uuid.UUID) error

	//User
	CheckIfUserOwnsRepo(userId uuid.UUID, repoId uuid.UUID) (bool, error)
	GetUserByID(id uuid.UUID) (*User, error)

	//Repository
	GetRepoIDByName(repoName string) (uuid.UUID, error)
	CheckIfRepoExists(repoName string) (bool, error)
	CheckIfRepoHasEventById(repoId uuid.UUID, eventTypeName uuid.UUID) (bool, error)

	//Installation
	GetInstallationIDByUser(userId uuid.UUID) (int, error)

	//Event
	CheckIfEventTypeExistsByName(eventType string) (bool, error)
	GetEventTypeByName(eventType string) (uuid.UUID, error)
	CreateEvent(event *Event) error
	DeleteEventForRepo(repoId uuid.UUID) error
}
type Notifications struct {
	UserId uuid.UUID `json:"user_id"`
	Repos  []struct {
		RepoName string   `json:"repo_name"`
		Events   []string `json:"events"`
	} `json:"repos"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	SoftDel   bool      `json:"soft_del"`
	CreatedAt time.Time `json:"created_at"`
}

type UserStore interface {
	GetUserByID(id uuid.UUID) (*User, error)
	CreateUser(user *User) error
	DeleteUser(id uuid.UUID) error
	GetUserByEmail(username string) (*User, error)
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

type InstallationStore interface {
	GetUserIdByUsername(username string) (uuid.UUID, error)
	CreateInstallation(userId uuid.UUID, installationId int) error
	CreateRepository(userId uuid.UUID, repoName string) error
	CheckIfRepoExists(repoName string) (bool, error)
	CheckIfInstallationExists(userId uuid.UUID) (bool, error)
	GetUserIdByInstallationId(installationId int) (uuid.UUID, error)
	RevokeUser(userId uuid.UUID) error
	GetInstallationIDByUser(userId uuid.UUID) (int, error)
	GetUserByID(id uuid.UUID) (*User, error)
	GetAllReposFromUserInNotificationSubscription(userId uuid.UUID) ([]*Repository, error)
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

type JwtToken struct {
	Token string `json:"token"`
}

type AccessToken struct {
	Token  string    `json:"token"`
	UserId uuid.UUID `json:"user_id"`
}

type EditNotification struct {
	RepoName string `json:"repo_name"`

	Events Events `json:"events"`

	UserID uuid.UUID `json:"user_id"`
}

type Events struct {
	Added  []string `json:"added"`
	Remove []string `json:"remove"`
}
