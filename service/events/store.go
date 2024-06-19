package events

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

func (s *Store) GetAllEvents() ([]types.EventType, error) {
	rows, err := s.db.Query("SELECT id, event_name, description_text FROM EventType ORDER BY event_name ASC;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]types.EventType, 0)

	for rows.Next() {
		e := types.EventType{}
		err = rows.Scan(&e.ID, &e.EventName, &e.Description)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (s *Store) GetAllEventsForRepo(reponame string) ([]types.EventType, error) {
	query := `
SELECT et.id, et.event_name, et.description_text
FROM Event e
JOIN Repository r ON e.repo_id = r.id
JOIN EventType et ON e.event_type = et.id
WHERE r.repo_name = ?
ORDER BY e.created_at ASC;
`

	rows, err := s.db.Query(query, reponame)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []types.EventType

	for rows.Next() {
		var event types.EventType
		if err := rows.Scan(&event.ID, &event.EventName, &event.Description); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Store) GetUserIDFromRepoName(reponame string) uuid.UUID {
	var userId string
	query := "SELECT user_id FROM Repository WHERE repo_name = ?;"
	err := s.db.QueryRow(query, reponame).Scan(&userId)

	if err != nil {
		return uuid.Nil
	}

	id, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil
	}

	return id
}

//todo

func (s *Store) GetNotificationSubscriptionId(userId uuid.UUID, reponame string) (uuid.UUID, error) {
	query := `
SELECT id
FROM NotificationSubscription
WHERE repo_id = (SELECT id FROM Repository WHERE repo_name = ? AND user_id = ?);

`

	var id string

	err := s.db.QueryRow(query, reponame, userId).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	subscriptionId, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}

	return subscriptionId, nil

}
