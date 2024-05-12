package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func CreateWebhook(username string, reponame string, events []string) error {

	godotenv.Load()
	// token := os.Getenv("GITHUB_API_KEY")

	// create a webhook
	serverUrl := os.Getenv("WEBHOOK_TEST")
	// githubUrl := "https://api.github.com/"
	// url := githubUrl + "repos/" + username + "/" + reponame + "/hooks"

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

	fmt.Println(string(payloadBytes))
	return err
}

func main() {
	CreateWebhook("username", "reponame", []string{"push"})
}
