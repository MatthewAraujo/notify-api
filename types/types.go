package types

import (
	"time"

	"github.com/google/uuid"
)

type GitHubError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
	Errors           []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

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
	Ref        string `json:"ref"`
	HookId     int    `json:"hook_id"`
	Repository struct {
		FullName string `json:"full_name"`
		Name     string `json:"name"`
		Owner    struct {
			Name  string `json:"login"`
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

type WelcomeEmail struct {
	Email      string
	Owner      string
	Repository string
}

type NotificationStore interface {
	// Notification
	CheckIfNotificationExists(id uuid.UUID) (bool, error)
	CheckIfNotificationExistsForUserId(userId uuid.UUID, repoId uuid.UUID) (bool, error)
	CreateNotification(notif *NotificationSubscription) error
	GetOwnerOfNotification(id uuid.UUID) (User, error)
	GetNotificationById(id uuid.UUID) (*NotificationSubscription, error)
	DeleteNotification(id uuid.UUID) error

	//User
	CheckIfUserOwnsRepo(userId uuid.UUID, repoId uuid.UUID) (bool, error)
	GetUserByID(id uuid.UUID) (*User, error)

	//Repository
	GetRepoIDByName(repoName string) (uuid.UUID, error)
	CheckIfRepoExists(repoName string) (bool, error)
	CheckIfRepoHasEventById(repoId uuid.UUID, eventTypeName uuid.UUID) (bool, error)
	GetRepoById(repoId uuid.UUID) (Repository, error)

	//Installation
	GetInstallationIDByUser(userId uuid.UUID) (int, error)

	//Event
	CheckIfEventTypeExistsByName(eventType string) (bool, error)
	GetEventTypeByName(eventType string) (uuid.UUID, error)
	CreateEvent(event *Event) error
	DeleteEventForRepo(repoId uuid.UUID) error

	//Webhook
	GetHookIdByRepoName(repoName string) (int, error)
}
type Notifications struct {
	UserId uuid.UUID `json:"user_id"`
	Repos  []struct {
		RepoName string   `json:"repo_name"`
		Events   []string `json:"events"`
	} `json:"repos"`
}

type User struct {
	ID        uuid.UUID
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	Email     string    `json:"email"`
	SoftDel   bool      `json:"soft_del"`
	CreatedAt time.Time `json:"created_at"`
}

type UserStore interface {
	GetUserByID(id uuid.UUID) (*User, error)
	CreateUser(user *User) error
	DeleteUser(id uuid.UUID) error
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
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
	GetUserIdByUsername(username string) (User, error)
	CreateInstallation(userId uuid.UUID, installationId int) error
	CreateRepository(userId uuid.UUID, repoName string) error
	CheckIfRepoExists(repoName string) (bool, error)
	CheckIfInstallationExists(userId uuid.UUID) (bool, error)
	GetUserIdByInstallationId(installationId int) (uuid.UUID, error)
	RevokeUser(userId uuid.UUID) error
	GetInstallationIDByUser(userId uuid.UUID) (int, error)
	GetUserByID(id uuid.UUID) (*User, error)
	GetAllReposFromUserInNotificationSubscription(userId uuid.UUID) ([]*Repository, error)

	//Webhook
	AddHookIdInNotificationSubscription(reponame string, hookId int) error
	CheckIfHookIdExistsInNotificationSubscription(hookId int) (bool, error)
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
	HookID    int       `json:"hook_id"`
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

type RepositoryStore interface {
	GetAllRepositoryForUser(username string) ([]ReposWithEvents, error)
}

type ReposWithEvents struct {
	RepoId   uuid.UUID   `json:"repo_id"`
	RepoName string      `json:"repo_name"`
	Events   []EventType `json:"events"`
}

type RepositoryEventType struct {
	RepoID      uuid.UUID
	RepoName    string
	EventTypeID uuid.UUID
	EventName   string
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
