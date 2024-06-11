package repository

import (
	"database/sql"

	"github.com/MatthewAraujo/notify/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// Repository
func (s *Store) GetAllRepositoryForUser(username string) ([]types.Repository, error) {
	query := `
 SELECT 
    r.id AS repo_id,
    r.repo_name,
    et.id AS event_type_id,
    et.event_name
FROM 
    Repository r
JOIN 
    User u ON r.user_id = u.id
JOIN 
    Event e ON e.repo_id = r.id
JOIN 
    EventType et ON e.event_type = et.id
WHERE 
    u.username = 'specific_username';`
	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []types.RepositoryEventType

	for rows.Next() {
		var ret types.RepositoryEventType
		if err := rows.Scan(&ret.RepoID, &ret.RepoName, &ret.EventTypeID, &ret.EventName); err != nil {
			return nil, err
		}
		results = append(results, ret)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	repos := make(map[string]*types.Repository)
	for _, ret := range results {
		if _, ok := repos[ret.RepoID]; !ok {
			repos[ret.RepoID] = &types.Repository{
				ID:       ret.RepoID,
				RepoName: ret.RepoName,
			}
		}
		repos[ret.RepoID].Events = append(repos[ret.RepoID].Events, types.EventType{
			ID:        ret.EventTypeID,
			EventName: ret.EventName,
		})

		var repositories []types.Repository
		for _, repo := range repos {
			repositories = append(repositories, *repo)
		}
		return repositories, nil

	}

func (s *Store) scanRowIntoRepository(rows *sql.Rows) (*types.Repository, error) {
	var r types.Repository
	err := rows.Scan(&r.ID, &r.RepoName)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
