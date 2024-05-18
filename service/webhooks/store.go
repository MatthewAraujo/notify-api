package webhooks

import (
	"database/sql"

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

func (s *Store) GetUserIdByUsername(username string) (uuid.UUID, error) {
	var userId uuid.UUID
	err := s.db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&userId)
	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}

func (s *Store) CreateInstallation(userId uuid.UUID, installationId int) error {
	_, err := s.db.Exec("INSERT INTO Installation (user_id, installation_id) VALUES (?, ?)", userId, installationId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CheckIfRepoExists(repoName string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM Repository WHERE repo_name = ?)", repoName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Store) CheckIfInstallationExists(userId uuid.UUID) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM Installation WHERE user_id = ?)", userId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Store) CreateRepository(userId uuid.UUID, repoName string) error {
	_, err := s.db.Exec("INSERT INTO Repository (id,user_id, repo_name) VALUES (?,?, ?)", uuid.New(), userId, repoName)
	if err != nil {
		return err
	}

	return nil
}
