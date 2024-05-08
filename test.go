package main

import (
	"encoding/json"
	"fmt"
)

type PushEvent struct {
	RepoName string `json:"repository"`
	Sender   string `json:"sender"`
	Commit   string `json:"head_commit"`
	Message  string `json:"message"`
	Email    string `json:"owner"`
}

func extractPushEventData(jsonData []byte) (PushEvent, error) {
	var pushEvent struct {
		Repo struct {
			Name string `json:"full_name"`
		} `json:"repository"`
		Sender struct {
			Login string `json:"login"`
		} `json:"sender"`
		HeadCommit struct {
			ID      string `json:"id"`
			Message string `json:"message"`
			Author  struct {
				Email string `json:"email"`
			} `json:"author"`
		} `json:"head_commit"`
	}

	err := json.Unmarshal(jsonData, &pushEvent)
	if err != nil {
		return PushEvent{}, err
	}

	event := PushEvent{
		RepoName: pushEvent.Repo.Name,
		Sender:   pushEvent.Sender.Login,
		Commit:   pushEvent.HeadCommit.ID,
		Message:  pushEvent.HeadCommit.Message,
		Email:    pushEvent.HeadCommit.Author.Email,
	}

	return event, nil
}

func main() {
	jsonData := []byte(`{"ref":"refs/heads/main","before":"96191d391e8f65e8cd2270f7ba587486956e74de","after":"86859e6ad72d831c7750ec75a6f19abdf33a3716","repository":{"id":797501575,"node_id":"R_kgDOL4johw","name":"test","full_name":"MatthewAraujo/test","private":true,"owner":{"name":"MatthewAraujo","email":"90223014+MatthewAraujo@users.noreply.github.com","login":"MatthewAraujo","id":90223014,"node_id":"MDQ6VXNlcjkwMjIzMDE0","avatar_url":"https://avatars.githubusercontent.com/u/90223014?v=4","gravatar_id":"","url":"https://api.github.com/users/MatthewAraujo","html_url":"https://github.com/MatthewAraujo","followers_url":"https://api.github.com/users/MatthewAraujo/followers","following_url":"https://api.github.com/users/MatthewAraujo/following{/other_user}","gists_url":"https://api.github.com/users/MatthewAraujo/gists{/gist_id}","starred_url":"https://api.github.com/users/MatthewAraujo/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/MatthewAraujo/subscriptions","organizations_url":"https://api.github.com/users/MatthewAraujo/orgs","repos_url":"https://api.github.com/users/MatthewAraujo/repos","events_url":"https://api.github.com/users/MatthewAraujo/events{/privacy}","received_events_url":"https://api.github.com/users/MatthewAraujo/received_events","type":"User","site_admin":false},"html_url":"https://github.com/MatthewAraujo/test","description":null,"fork":false,"url":"https://github.com/MatthewAraujo/test","forks_url":"https://api.github.com/repos/MatthewAraujo/test/forks","keys_url":"https://api.github.com/repos/MatthewAraujo/test/keys{/key_id}","collaborators_url":"https://api.github.com/repos/MatthewAraujo/test/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/MatthewAraujo/test/teams","hooks_url":"https://api.github.com/repos/MatthewAraujo/test/hooks","issue_events_url":"https://api.github.com/repos/MatthewAraujo/test/issues/events{/number}","events_url":"https://api.github.com/repos/MatthewAraujo/test/events","assignees_url":"https://api.github.com/repos/MatthewAraujo/test/assignees{/user}","branches_url":"https://api.github.com/repos/MatthewAraujo/test/branches{/branch}","tags_url":"https://api.github.com/repos/MatthewAraujo/test/tags","blobs_url":"https://api.github.com/repos/MatthewAraujo/test/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/MatthewAraujo/test/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/MatthewAraujo/test/git/refs{/sha}","trees_url":"https://api.github.com/repos/MatthewAraujo/test/git/trees{/sha}","statuses_url":"https://api.github.com/repos/MatthewAraujo/test/statuses/{sha}","languages_url":"https://api.github.com/repos/MatthewAraujo/test/languages","stargazers_url":"https://api.github.com/repos/MatthewAraujo/test/stargazers","contributors_url":"https://api.github.com/repos/MatthewAraujo/test/contributors","subscribers_url":"https://api.github.com/repos/MatthewAraujo/test/subscribers","subscription_url":"https://api.github.com/repos/MatthewAraujo/test/subscription","commits_url":"https://api.github.com/repos/MatthewAraujo/test/commits{/sha}","git_commits_url":"https://api.github.com/repos/MatthewAraujo/test/git/commits{/sha}","comments_url":"https://api.github.com/repos/MatthewAraujo/test/comments{/number}","issue_comment_url":"https://api.github.com/repos/MatthewAraujo/test/issues/comments{/number}","contents_url":"https://api.github.com/repos/MatthewAraujo/test/contents/{+path}","compare_url":"https://api.github.com/repos/MatthewAraujo/test/compare/{base}...{head}","merges_url":"https://api.github.com/repos/MatthewAraujo/test/merges","archive_url":"https://api.github.com/repos/MatthewAraujo/test/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/MatthewAraujo/test/downloads","issues_url":"https://api.github.com/repos/MatthewAraujo/test/issues{/number}","pulls_url":"https://api.github.com/repos/MatthewAraujo/test/pulls{/number}","milestones_url":"https://api.github.com/repos/MatthewAraujo/test/milestones{/number}","notifications_url":"https://api.github.com/repos/MatthewAraujo/test/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/MatthewAraujo/test/labels{/name}","releases_url":"https://api.github.com/repos/MatthewAraujo/test/releases{/id}","deployments_url":"https://api.github.com/repos/MatthewAraujo/test/deployments","created_at":1715129890,"updated_at":"2024-05-08T01:06:40Z","pushed_at":1715130441,"git_url":"git://github.com/MatthewAraujo/test.git","ssh_url":"git@github.com:MatthewAraujo/test.git","clone_url":"https://github.com/MatthewAraujo/test.git","svn_url":"https://github.com/MatthewAraujo/test","homepage":null,"size":0,"stargazers_count":0,"watchers_count":0,"language":"Go","has_issues":true,"has_projects":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"has_discussions":false,"forks_count":0,"mirror_url":null,"archived":false,"disabled":false,"open_issues_count":0,"license":null,"allow_forking":true,"is_template":false,"web_commit_signoff_required":false,"topics":[],"visibility":"private","forks":0,"open_issues":0,"watchers":0,"default_branch":"main","stargazers":0,"master_branch":"main"},"pusher":{"name":"MatthewAraujo","email":"90223014+MatthewAraujo@users.noreply.github.com"},"sender":{"login":"MatthewAraujo","id":90223014,"node_id":"MDQ6VXNlcjkwMjIzMDE0","avatar_url":"https://avatars.githubusercontent.com/u/90223014?v=4","gravatar_id":"","url":"https://api.github.com/users/MatthewAraujo","html_url":"https://github.com/MatthewAraujo","followers_url":"https://api.github.com/users/MatthewAraujo/followers","following_url":"https://api.github.com/users/MatthewAraujo/following{/other_user}","gists_url":"https://api.github.com/users/MatthewAraujo/gists{/gist_id}","starred_url":"https://api.github.com/users/MatthewAraujo/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/MatthewAraujo/subscriptions","organizations_url":"https://api.github.com/users/MatthewAraujo/orgs","repos_url":"https://api.github.com/users/MatthewAraujo/repos","events_url":"https://api.github.com/users/MatthewAraujo/events{/privacy}","received_events_url":"https://api.github.com/users/MatthewAraujo/received_events","type":"User","site_admin":false},"created":false,"deleted":false,"forced":false,"base_ref":null,"compare":"https://github.com/MatthewAraujo/test/compare/96191d391e8f...86859e6ad72d","commits":[{"id":"bc80efdc9c83aadfd2b7ab4874f92f7387d37847","tree_id":"b9270df7070cc6a5e7dbdec610a7ce4f54c47b20","distinct":false,"message":"removing","timestamp":"2024-05-07T22:07:06-03:00","url":"https://github.com/MatthewAraujo/test/commit/bc80efdc9c83aadfd2b7ab4874f92f7387d37847","author":{"name":"MatthewAraujo","email":"matthewaraujo20@gmail.com","username":"MatthewAraujo"},"committer":{"name":"MatthewAraujo","email":"matthewaraujo20@gmail.com","username":"MatthewAraujo"},"added":[],"removed":[],"modified":["main.go"]},{"id":"86859e6ad72d831c7750ec75a6f19abdf33a3716","tree_id":"b9270df7070cc6a5e7dbdec610a7ce4f54c47b20","distinct":true,"message":"Merge pull request #1 from MatthewAraujo/test\n\nremoving","timestamp":"2024-05-07T22:07:21-03:00","url":"https://github.com/MatthewAraujo/test/commit/86859e6ad72d831c7750ec75a6f19abdf33a3716","author":{"name":"Matthew Araujo","email":"90223014+MatthewAraujo@users.noreply.github.com","username":"MatthewAraujo"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["main.go"]}],"head_commit":{"id":"86859e6ad72d831c7750ec75a6f19abdf33a3716","tree_id":"b9270df7070cc6a5e7dbdec610a7ce4f54c47b20","distinct":true,"message":"Merge pull request #1 from MatthewAraujo/test\n\nremoving","timestamp":"2024-05-07T22:07:21-03:00","url":"https://github.com/MatthewAraujo/test/commit/86859e6ad72d831c7750ec75a6f19abdf33a3716","author":{"name":"Matthew Araujo","email":"90223014+MatthewAraujo@users.noreply.github.com","username":"MatthewAraujo"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["main.go"]}}`)

	event, err := extractPushEventData(jsonData)
	if err != nil {
		fmt.Println("Erro ao extrair dados do evento:", err)
		return
	}

	fmt.Println("Repositório:", event.RepoName)
	fmt.Println("Enviado por:", event.Sender)
	fmt.Println("Commit:", event.Commit)
	fmt.Println("Mensagem do commit:", event.Message)
	fmt.Println("Email do dono do projeto:", event.Email)
}
