package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MatthewAraujo/notify/config"
	"github.com/MatthewAraujo/notify/db"
	"github.com/MatthewAraujo/notify/types"
	"github.com/google/uuid"
)

func RequestAccessToken(userId uuid.UUID, installationID string, jwtToken string) (string, error) {

	accessToken, err := GetAccessToken(userId)
	if err != nil {
		if err.Error() == "access token not found" {
			log.Printf("Access token not found in database")
		} else {
			if accessToken.Token != "" {
				return accessToken.Token, nil
			}
		}
	}
	// Create HTTP client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installationID), nil)
	if err != nil {
		return "", err
	}

	// Add headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Parse the response
	var accessTokenResp struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&accessTokenResp)
	if err != nil {
		return "", err
	}

	return accessTokenResp.Token, nil
}

func GetAccessToken(id uuid.UUID) (*types.AccessToken, error) {
	db, err := db.NewMySQLStorage(config.Envs.TursoURl)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT token FROM AccessToken WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	at := new(types.AccessToken)

	for rows.Next() {
		at, err = scanRowIntoAccessTokenWithUserID(rows)
		if err != nil {
			return nil, err
		}
	}

	if at.Token == "" {
		return nil, fmt.Errorf("access token not found")
	}

	return at, nil
}

func scanRowIntoAccessTokenWithUserID(rows *sql.Rows) (*types.AccessToken, error) {
	at := new(types.AccessToken)
	u := new(types.User)

	if err := rows.Scan(&at.Token, &u.ID); err != nil {
		return nil, err
	}

	return at, nil
}
