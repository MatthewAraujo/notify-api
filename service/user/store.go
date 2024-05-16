package user

import (
	"database/sql"
	"fmt"
	"time"

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

func (s *Store) GetUserByEmail(username string) (*types.User, error) {
	rows, err := s.db.Query("SELECT id, email FROM user WHERE username = $1", username)
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

func (s *Store) CreateUser(user *types.User) error {
	_, err := s.db.Exec("INSERT INTO user (id, username, email, created_at) VALUES ($1, $2, $3, $4)", user.ID, user.Username, user.Email, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByID(id uuid.UUID) (*types.User, error) {
	rows, err := s.db.Query("SELECT id, username FROM user WHERE id = $1", id)
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
func (s *Store) DeleteUser(id uuid.UUID) error {
	_, err := s.db.Exec("UPDATE user SET deleted_at = $1 WHERE id = $2", time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	var user types.User
	if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}
