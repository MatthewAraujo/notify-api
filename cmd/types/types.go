package types

import "github.com/google/uuid"

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

type NotificationStore interface {
	GetUserByID(id uuid.UUID) (Notifications, error)
	GetRepositoryByUserID(id uuid.UUID, reponame string) (Notifications, error)
	CreateNotification(notif Notifications) error
}
type Notifications struct {
	UserId   uuid.UUID `json:"user_id"`
	RepoName string    `json:"repository_name"`
	Events   Events    `json:"events"`
}

type Events struct {
	BranchCreation                bool
	TagCreation                   bool
	BranchDeletion                bool
	TagDeletion                   bool
	BranchProtectionConfig        bool
	BranchProtectionRules         bool
	CheckRuns                     bool
	CheckSuites                   bool
	CodeScanningAlerts            bool
	CollaboratorChange            bool
	CommitComments                bool
	DependabotAlerts              bool
	DeployKeys                    bool
	DeploymentStatuses            bool
	Deployments                   bool
	DiscussionComments            bool
	Discussions                   bool
	Forks                         bool
	IssueComments                 bool
	Issues                        bool
	Labels                        bool
	MergeGroups                   bool
	Meta                          bool
	Milestones                    bool
	Packages                      bool
	PageBuilds                    bool
	ProjectCards                  bool
	ProjectColumns                bool
	Projects                      bool
	PullRequestReviewComments     bool
	PullRequestReviewThreads      bool
	PullRequestReviews            bool
	PullRequests                  bool
	Pushes                        bool
	RegistryPackages              bool
	Releases                      bool
	Repositories                  bool
	RepositoryAdvisories          bool
	RepositoryImports             bool
	RepositoryRulesets            bool
	RepositoryVulnerabilityAlerts bool
	SecretScanningAlertLocations  bool
	SecretScanningAlerts          bool
	SecurityAndAnalyses           bool
	Stars                         bool
	Statuses                      bool
	TeamAdds                      bool
	VisibilityChanges             bool
	Watches                       bool
	Wiki                          bool
	WorkflowJobs                  bool
	WorkflowRuns                  bool
}
