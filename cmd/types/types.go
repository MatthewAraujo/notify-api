package types

type GithubWebhooks struct {
	RepoName string `json:"full_name"`
	Sender   string `json:"login"`
	Commit   string `json:"message"`
}

type SendEmail struct {
	RepoName string
	Sender   string
	Commit   string
	Email    string
}
