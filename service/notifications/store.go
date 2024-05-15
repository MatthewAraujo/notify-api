package notifications

import (
	"database/sql"

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

type User struct {
	ID   int
	Name string
}

func (s *Store) CreateNotification(notif *types.Notifications) error {
	return nil
}

func (s *Store) GetRepositoryByUserID(id uuid.UUID, reponame string) (*types.Notifications, error) {
	return nil, nil
}

func (s *Store) GetUserByID(id uuid.UUID) (*types.Notifications, error) {
	return nil, nil
}

func (s *Store) scanRowIntoUser(rows *sql.Rows) (*User, error) {
	var user User
	if err := rows.Scan(&user.ID, &user.Name); err != nil {
		return nil, err
	}
	return &user, nil
}
