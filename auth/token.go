package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MatthewAraujo/notify/config"
	"github.com/MatthewAraujo/notify/db"
	"github.com/MatthewAraujo/notify/encrypt"
	"github.com/MatthewAraujo/notify/types"
	"github.com/google/uuid"
)

func RequestAccessToken(userId uuid.UUID, installationID int, jwtToken string) (string, error) {

	accessToken, err := GetAccessToken(userId)
	if err != nil {
		log.Printf("Error getting access token: %s", err)
		if err.Error() == "access token not found" {
			log.Printf("Access token not found in database")
		}
	} else {
		if accessToken != "" {
			token, err := encrypt.DecryptToken(accessToken)
			if err != nil {
				return "", err
			}
			return token, nil
		}
	}

	log.Printf("Access token not found in database, requesting new token")
	// Create request
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID), bytes.NewBuffer([]byte{}))
	if err != nil {
		return "", err
	}

	// Add headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY4NTk3NTksImlhdCI6MTcxNjg1OTE1OSwiaXNzIjoiODk3MzkyIn0.JCf3DZr7FyZ7X15noD0eX4q2ih1rD0CTb_RMLkF_yLDcsfA-QR3iwZXaJqrgbWpDojXaQPGtD7_Ki5In1LMTrjQAYywYYlEm3DOmR03rai17au5cxOWIatS_sHHVG6fUg8FFL0k_75hnAGxT1nZiALeaQd0H3RLihe7c3WI5CZCkKoCna0y4U-hapHNzuvn1mxaWNume1SbOEqGj_ImaoIy2ZE0puALN6n4_dTgUkLB7N961AL4Io-KaglkVfy0Icinkd8oxqBxklsK0U6kXsjMT0VLmlua44Jl1X3Ha9NOhky1TuqIsPQ5lESb2PYh0TQtLmz3keLYQ0tDvwUhqTQ")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Create HTTP client
	client := &http.Client{}
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

	// Insert the access token into the database
	token, err := encrypt.EncryptToken(accessTokenResp.Token)
	if err != nil {
		return "", err
	}

	err = insertAccessToken(&types.AccessToken{
		Token:  token,
		UserId: userId,
	})

	if err != nil {
		return "", err
	}

	return accessTokenResp.Token, nil
}

func GetAccessToken(id uuid.UUID) (string, error) {
	db, err := db.NewMySQLStorage(config.Envs.TursoURl)
	if err != nil {
		return "", err
	}

	var token string
	err = db.QueryRow("SELECT token FROM AccessToken WHERE user_id = ?", id).Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("access token not found")
		}
		return "", err
	}

	if token == "" {
		return "", fmt.Errorf("access token not found")
	}

	return token, nil
}

func insertAccessToken(at *types.AccessToken) error {
	db, err := db.NewMySQLStorage(config.Envs.TursoURl)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO AccessToken (token, user_id) VALUES (?, ?)", at.Token, at.UserId)
	if err != nil {
		return err
	}

	return nil
}
