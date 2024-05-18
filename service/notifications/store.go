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

func (s *Store) CreateNotification(notification *types.NotificationSubscription) error {
	_, err := s.db.Exec("INSERT INTO NotificationSubscription (id,repo_id) VALUES (?, ?)", uuid.New(), notification.RepoID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateEvent(event *types.Event) error {
	_, err := s.db.Exec("INSERT INTO event (id, repo_id,event_type) VALUES (?, ?, ?)", uuid.New(), event.RepoID, event.EventType)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetEventTypeByName(name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.db.QueryRow("SELECT id FROM EventType WHERE event_name = ?", name).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *Store) GetInstallationIDByUser(id uuid.UUID) (int, error) {
	var installationID int
	err := s.db.QueryRow("SELECT installation_id FROM installation WHERE user_id = ?", id).Scan(&installationID)
	if err != nil {
		return 0, err
	}

	return installationID, nil
}

func (s *Store) GetRepoIDByName(name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.db.QueryRow("SELECT id FROM repository WHERE repo_name = ?", name).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("repo not found")
	}

	return id, nil
}

func (s *Store) CheckIfNotificationExists(userID uuid.UUID, repoID uuid.UUID) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM NotificationSubscription WHERE repo_id = ?)", repoID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Store) CheckIfRepoExists(name string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM repository WHERE repo_name = ?)", name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Store) scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	var user types.User
	if err := rows.Scan(&user.ID, &user.Username); err != nil {
		return nil, err
	}
	return &user, nil
}
