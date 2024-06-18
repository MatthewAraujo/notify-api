package repository

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

// Repository
func (s *Store) GetAllRepositoryForUser(username string) ([]types.ReposWithEvents, error) {
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
    u.username = ? ;`
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

	repoMap := make(map[uuid.UUID]*types.ReposWithEvents)
	for _, ret := range results {
		if _, ok := repoMap[ret.RepoID]; !ok {
			repoMap[ret.RepoID] = &types.ReposWithEvents{
				RepoId:   ret.RepoID,
				RepoName: ret.RepoName,
				Events:   []types.EventType{},
			}
		}
		repoMap[ret.RepoID].Events = append(repoMap[ret.RepoID].Events, types.EventType{
			ID:        ret.EventTypeID,
			EventName: ret.EventName,
		})
	}

	var repos []types.ReposWithEvents
	for _, v := range repoMap {
		repos = append(repos, *v)
	}

	return repos, nil
}

func (s *Store) GetAllReposForUser(username string) ([]types.Repository, error) {
	query := `SELECT r.id, r.repo_name FROM Repository r JOIN User u ON r.user_id = u.id WHERE u.username = ?;`
	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []types.Repository
	for rows.Next() {
		repo, err := s.scanRowIntoRepository(rows)
		if err != nil {
			return nil, err
		}
		repos = append(repos, *repo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}

func (s *Store) IsRepoSubscribed(username string, repoId uuid.UUID) (bool, error) {
	query := `
SELECT COUNT(*)
FROM NotificationSubscription ns
JOIN Repository r ON ns.repo_id = r.id
JOIN User u ON r.user_id = u.id
WHERE u.username = ? AND r.id = ? AND ns.removed = FALSE;
`
	var count int
	err := s.db.QueryRow(query, username, repoId).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *Store) scanRowIntoRepository(rows *sql.Rows) (*types.Repository, error) {
	var r types.Repository
	err := rows.Scan(&r.ID, &r.RepoName)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
