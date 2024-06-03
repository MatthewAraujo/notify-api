package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/MatthewAraujo/notify/auth"
	"github.com/MatthewAraujo/notify/config"
	"github.com/MatthewAraujo/notify/types"
	"github.com/google/uuid"
)

func CreateWebhook(installationId int, username string, userId uuid.UUID, reponame string, events []string) error {

	token, err := generateAccessToken(installationId, userId)
	if err != nil {
		return err
	}

	// create a webhook
	serverUrl := config.Envs.WebhookUrl
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

	return nil
}

func UpdateWebhook(username string, userID uuid.UUID, reponame string, events types.Events, db types.NotificationStore) error {
	// get the installation id
	installationId, err := db.GetInstallationIDByUser(userID)
	if err != nil {
		return err
	}

	// get the hook id
	hookId, err := db.GetHookIdByRepoName(reponame)
	if err != nil {
		return err
	}

	addedEvents := events.Added
	removedEvents := events.Remove

	err = updateWebhook(installationId, username, userID, reponame, addedEvents, removedEvents, hookId)
	if err != nil {
		return err
	}

	return nil
}

func updateWebhook(installationId int, username string, userId uuid.UUID, reponame string, addedEvents, removedEvents []string, hookId int) error {
	token, err := generateAccessToken(installationId, userId)
	if err != nil {
		return err
	}

	// create a webhook
	githubUrl := "https://api.github.com/"
	url := githubUrl + "repos/" + username + "/" + reponame + "/hooks" + fmt.Sprintf("/%d", hookId)

	log.Printf("URL: %s", url)
	log.Printf("access token: %s", token)
	err = updatePayloadToGithub(url, token, addedEvents, removedEvents)
	if err != nil {
		return err
	}

	return nil
}

func DeleteWebhook(userId uuid.UUID, installationId int, username string, reponame string, hookId int) error {
	token, err := generateAccessToken(installationId, userId)
	if err != nil {
		return err
	}

	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://api.github.com/repos/%s/%s/hooks/%d", username, reponame, hookId), nil)
	if err != nil {
		return err
	}

	// Add headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusNoContent {
		var message struct {
			Message string `json:"message"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
			return fmt.Errorf("unexpected response status: %s", resp.Status)
		}
		return fmt.Errorf("GitHub API error: %s", message.Message)

	}

	return nil
}

func DeleteAllWebhooks(userId uuid.UUID, db types.InstallationStore) error {
	log.Printf("Deleting webhooks for user %s", userId)
	// get the installation id
	installationId, err := db.GetInstallationIDByUser(userId)
	if err != nil {
		return err
	}

	// get the username
	user, err := db.GetUserByID(userId)
	if err != nil {
		return err
	}

	// get all repos for the user that is on the NotificationSubscription
	repos, err := db.GetAllReposFromUserInNotificationSubscription(userId)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		log.Printf("Deleting webhook for %s", repo.RepoName)
		err = deleteWebhook(installationId, user, repo.RepoName, userId)
		if err != nil {
			return err

		}
	}
	return nil
}

func deleteWebhook(installationId int, user *types.User, reponame string, userId uuid.UUID) error {

	token, err := generateAccessToken(installationId, userId)
	if err != nil {
		return err
	}

	// create a webhook
	githubUrl := "https://api.github.com/"
	url := githubUrl + "repos/" + user.Username + "/" + reponame + "/hooks"

	err = deletePayloadToGitHub(url, token)
	if err != nil {
		return err

	}

	return nil
}

func deletePayloadToGitHub(url, token string) error {

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	// Set the headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var ghErr types.GitHubError
		if err := json.Unmarshal(body, &ghErr); err != nil {
			return fmt.Errorf("unexpected response status: %s", resp.Status)
		}
		return fmt.Errorf("GitHub API error: %s, Details: %+v", ghErr.Message, ghErr.Errors)

	}

	return nil
}

// generate access token
func generateAccessToken(installationId int, userId uuid.UUID) (string, error) {
	jwt, err := auth.GenerateJWT()
	if err != nil {
		return "", err
	}

	accessToken, err := auth.RequestAccessToken(userId, installationId, jwt)
	if err != nil {
		return "", err
	}

	return accessToken, nil
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

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		var message struct {
			Message string `json:"message"`
		}

		if err := json.Unmarshal(body, &message); err != nil {
			return fmt.Errorf("unexpected response status: %s", resp.Status)
		}

		return fmt.Errorf("GitHub API error: %s", message.Message)
	}

	return nil
}

func updatePayloadToGithub(url, token string, addedEvents, removedEvents []string) error {
	payload := map[string]interface{}{
		"active":        true,
		"add_events":    addedEvents,
		"remove_events": removedEvents,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}
	// Set the headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode)
		var messages struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(body, &messages); err != nil {
			return fmt.Errorf("unexpected response status: %s", resp.Status)
		}

		return fmt.Errorf("GitHub API error: %s", messages.Message)
	}

	return nil
}
