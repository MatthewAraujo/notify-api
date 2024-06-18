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

func (s *Store) GetAllEventsForRepo(repoId uuid.UUID) ([]types.EventType, error) {
	query := `
SELECT * FROM Event WHERE repo_id = ? ORDER BY event_name ASC
;
`

	rows, err := s.db.Query(query, repoId)
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
