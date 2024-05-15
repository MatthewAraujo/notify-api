package notifications

import (
	"database/sql"
	"fmt"
	"os"

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

func (s *Store) QueryUsers(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		users = append(users, user)
		fmt.Println(user.ID, user.Name)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}
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
