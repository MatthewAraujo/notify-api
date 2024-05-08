package types

type GithubWebhooks struct {
	Repository struct {
		FullName string `json:"full_name"`
		Owner    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"owner"`
	} `json:"repository"`
	Commits []struct {
		Message string `json:"message"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			// Outros campos do autor do commit, se necessário
		} `json:"author"`
	} `json:"commits"`
}

type SendEmail struct {
	RepoName string
	Sender   string
	Commit   string
	Email    string
}
