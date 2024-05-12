package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func CreateWebhook(username string, reponame string, events []string) error {

	godotenv.Load()
	token := os.Getenv("GITHUB_API_KEY")

	// create a webhook
	serverUrl := os.Getenv("WEBHOOK_TEST")
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
