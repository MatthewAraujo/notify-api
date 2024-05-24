package webhooks

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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

func (s *Store) GetUserIdByInstallationId(installationId int) (uuid.UUID, error) {
	log.Printf("Getting user id by installation id %d", installationId)
	var userId uuid.UUID
	err := s.db.QueryRow("SELECT user_id FROM Installation WHERE installation_id = ?", installationId).Scan(&userId)
	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}

func (s *Store) RevokeUser(userId uuid.UUID) error {

	exists, err := s.CheckIfInstallationExists(userId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("instalação não encontrada para o usuário")
	}

	if err := s.RevokeInstallation(userId); err != nil {
		return err
	}

	if err := s.DeleteRepositoriesByUserId(userId); err != nil {
		return err
	}

	if err := s.RemoveNotificationSubscriptionsByUserId(userId); err != nil {
		return err
	}

	if err := s.DeleteAccessTokensByUserId(userId); err != nil {
		return err
	}

	return nil
}

func (s *Store) RevokeInstallation(userId uuid.UUID) error {
	_, err := s.db.Exec("UPDATE Installation SET revoked_at = ? WHERE user_id = ?", time.Now(), userId)
	return err
}

func (s *Store) DeleteRepositoriesByUserId(userId uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM Repository WHERE user_id = ?", userId)
	return err
}

func (s *Store) RemoveNotificationSubscriptionsByUserId(userId uuid.UUID) error {
	_, err := s.db.Exec(`UPDATE NotificationSubscription SET removed = TRUE, updated_at = ? 
		WHERE repo_id IN (SELECT id FROM Repository WHERE user_id = ?)`, time.Now(), userId)
	return err
}

func (s *Store) DeleteAccessTokensByUserId(userId uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM AccessToken WHERE user_id = ?", userId)
	return err
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
func (s *Store) GetAllReposFromUserInNotificationSubscription(userID uuid.UUID) ([]*types.Repository, error) {
	query := `
		SELECT r.id, r.repo_name, r.user_id
		FROM NotificationSubscription ns
		JOIN Repository r ON ns.repo_id = r.id
		WHERE r.user_id = ?
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []*types.Repository

	for rows.Next() {
		var repo types.Repository
		if err := rows.Scan(&repo.ID, &repo.RepoName, &repo.UserID); err != nil {
			return nil, err
		}
		repos = append(repos, &repo)
	}

	return repos, nil
}

func (s *Store) GetInstallationIDByUser(id uuid.UUID) (int, error) {
	var installationID int
	err := s.db.QueryRow("SELECT installation_id FROM installation WHERE user_id = ?", id).Scan(&installationID)
	if err != nil {
		return 0, err
	}

	return installationID, nil
}

func (s *Store) CheckIfHookIdExistsInNotificationSubscription(hookId int) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM NotificationSubscription WHERE hook_id = ?)", hookId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Store) AddHookIdInNotificationSubscription(reponame string, hookId int) error {
	query := `
		UPDATE NotificationSubscription
		SET hook_id = ?
		WHERE repo_id = (SELECT id FROM Repository WHERE repo_name = ?)
	`

	_, err := s.db.Exec(query, hookId, reponame)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("repo not found")
		}

		return err
	}

	return err

}

func (s *Store) scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	var user types.User
	if err := rows.Scan(&user.ID, &user.Username); err != nil {
		return nil, err
	}
	return &user, nil
}
