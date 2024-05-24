package notifications

import (
	"database/sql"
	"fmt"

	"github.com/MatthewAraujo/notify/types"
	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// Notification
func (s *Store) GetOwnerOfNotification(id uuid.UUID) (uuid.UUID, error) {
	var userID uuid.UUID
	// select repo_id from notification where id = ?
	// select user_id from repository where id = ?

	err := s.db.QueryRow("SELECT user_id FROM repository WHERE id = (SELECT repo_id FROM NotificationSubscription WHERE id = ?)", id).Scan(&userID)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

func (s *Store) DeleteNotification(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM NotificationSubscription WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil

}

func (s *Store) CreateNotification(notification *types.NotificationSubscription) error {
	_, err := s.db.Exec("INSERT INTO NotificationSubscription (id,repo_id) VALUES (?, ?)", uuid.New(), notification.RepoID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CheckIfNotificationExists(id uuid.UUID) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM NotificationSubscription WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Store) CheckIfNotificationExistsForUserId(userID uuid.UUID, repoID uuid.UUID) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM NotificationSubscription WHERE repo_id = ?)", repoID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// User

func (s *Store) CheckIfUserOwnsRepo(userID uuid.UUID, repoID uuid.UUID) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM repository WHERE user_id = ? AND id = ?)", userID, repoID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Store) GetUserByID(id uuid.UUID) (*types.User, error) {
	rows, err := s.db.Query("SELECT id, username FROM user WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u := new(types.User)

	for rows.Next() {
		u, err = s.scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == uuid.Nil {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

// Event
func (s *Store) CreateEvent(event *types.Event) error {
	_, err := s.db.Exec("INSERT INTO event (id, repo_id,event_type) VALUES (?, ?, ?)", uuid.New(), event.RepoID, event.EventType)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteEventForRepo(repoID uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM event WHERE repo_id = ?", repoID)
	if err != nil {
		return err
	}
	return nil
}

// EventType
func (s *Store) GetEventTypeByName(name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.db.QueryRow("SELECT id FROM EventType WHERE event_name = ?", name).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *Store) CheckIfEventTypeExistsByName(name string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM EventType WHERE event_name = ?)", name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Installation
func (s *Store) GetInstallationIDByUser(id uuid.UUID) (int, error) {
	var installationID int
	err := s.db.QueryRow("SELECT installation_id FROM installation WHERE user_id = ?", id).Scan(&installationID)
	if err != nil {
		return 0, err
	}

	return installationID, nil
}

// Repo

func (s *Store) GetRepoIDByName(name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.db.QueryRow("SELECT id FROM repository WHERE repo_name = ?", name).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("repo not found")
	}

	return id, nil
}

func (s *Store) CheckIfRepoExists(name string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM repository WHERE repo_name = ?)", name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
func (s *Store) CheckIfRepoHasEventById(repoID uuid.UUID, eventName uuid.UUID) (bool, error) {
	var exists bool

	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM event WHERE repo_id = ? AND event_type = ?)", repoID, eventName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Webhook
func (s *Store) GetHookIdByRepoName(name string) (int, error) {

	query := `
		SELECT hook_id
		FROM NotificationSubscription
		WHERE repo_id = (SELECT id FROM Repository WHERE repo_name = ?)`
	var id int
	err := s.db.QueryRow(query, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	var user types.User
	if err := rows.Scan(&user.ID, &user.Username); err != nil {
		return nil, err
	}
	return &user, nil
}
