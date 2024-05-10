package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	serverUrl    string
	contentType  string
	insecure_ssl string
}

func CreateWebhook(username string, reponame string, events []string) error {

	godotenv.Load()
	token := os.Getenv("GITHUB_API_KEY")

	// create a webhook
	serverUrl := os.Getenv("WEBHOOK_TEST")
	githubUrl := "https://api.github.com/"
	url := githubUrl + "repos/" + username + "/" + reponame + "/hooks"

	config := Config{
		serverUrl:    serverUrl,
		contentType:  "json",
		insecure_ssl: "0",
	}

	payload := map[string]interface{}{
		"name":   "web",
		"active": true,
		"events": events,
		"config": map[string]interface{}{
			"url":          config.serverUrl,
			"content_type": config.contentType,
			"insecure_ssl": config.insecure_ssl,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// send the payload to the github api
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

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
