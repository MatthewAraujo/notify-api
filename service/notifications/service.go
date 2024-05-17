package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MatthewAraujo/notify/jwt"
	"github.com/joho/godotenv"
)

func CreateWebhook(username string, reponame string, events []string) error {

	godotenv.Load()
	token, err := generateAccessToken()
	if err != nil {
		return err
	}

	// create a webhook
	serverUrl := "https://scarce-joystick-04.webhook.cool"
	githubUrl := "https://api.github.com/"
	url := githubUrl + "repos/" + username + "/" + reponame + "/hooks"

	payload := map[string]interface{}{
		"name":   "web",
		"active": true,
		"events": events,
		"config": map[string]interface{}{
			"url":          serverUrl,
			"content_type": "json",
			"insecure_ssl": 0,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = sendPayloadToGitHub(url, token, payloadBytes)
	if err != nil {
		return err

	}

	return err
}

// generate access token
func generateAccessToken() (string, error) {
	jwt, err := jwt.GenerateJWT()
	if err != nil {
		return "", err
	}

	// how i will get this????
	// i wuill get this from the database
	installationID := "50730929"
	accessToken, err := getInstallationAccessToken(installationID, jwt)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func getInstallationAccessToken(installationID, jwtToken string) (string, error) {
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

func sendPayloadToGitHub(url, token string, payloadBytes []byte) error {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}

	// Add headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
