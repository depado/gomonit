package models

import "time"

// GHUser represents a github user
type GHUser struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

// Commit represents a single commit
type Commit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
		Tree    struct {
			Sha string `json:"sha"`
			URL string `json:"url"`
		} `json:"tree"`
		URL          string `json:"url"`
		CommentCount int    `json:"comment_count"`
	} `json:"commit"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	CommentsURL string `json:"comments_url"`
	Author      GHUser `json:"author"`
	Commiter    GHUser `json:"committer"`
	Parents     []struct {
		Sha     string `json:"sha"`
		URL     string `json:"url"`
		HTMLURL string `json:"html_url"`
	} `json:"parents"`
}

// Commits represents multiple commits
type Commits []Commit

// GHRepo represents a single Github repository
type GHRepo struct {
	ID               int            `json:"id"`
	Owner            GHUser         `json:"owner"`
	Name             string         `json:"name"`
	FullName         string         `json:"full_name"`
	Description      string         `json:"description"`
	Private          bool           `json:"private"`
	Fork             bool           `json:"fork"`
	URL              string         `json:"url"`
	HTMLURL          string         `json:"html_url"`
	ArchiveURL       string         `json:"archive_url"`
	AssigneesURL     string         `json:"assignees_url"`
	BlobsURL         string         `json:"blobs_url"`
	BranchesURL      string         `json:"branches_url"`
	CloneURL         string         `json:"clone_url"`
	CollaboratorsURL string         `json:"collaborators_url"`
	CommentsURL      string         `json:"comments_url"`
	CommitsURL       string         `json:"commits_url"`
	CompareURL       string         `json:"compare_url"`
	ContentsURL      string         `json:"contents_url"`
	ContributorsURL  string         `json:"contributors_url"`
	DeploymentsURL   string         `json:"deployments_url"`
	DownloadsURL     string         `json:"downloads_url"`
	EventsURL        string         `json:"events_url"`
	ForksURL         string         `json:"forks_url"`
	GitCommitsURL    string         `json:"git_commits_url"`
	GitRefsURL       string         `json:"git_refs_url"`
	GitTagsURL       string         `json:"git_tags_url"`
	GitURL           string         `json:"git_url"`
	HooksURL         string         `json:"hooks_url"`
	IssueCommentURL  string         `json:"issue_comment_url"`
	IssueEventsURL   string         `json:"issue_events_url"`
	IssuesURL        string         `json:"issues_url"`
	KeysURL          string         `json:"keys_url"`
	LabelsURL        string         `json:"labels_url"`
	LanguagesURL     string         `json:"languages_url"`
	MergesURL        string         `json:"merges_url"`
	MilestonesURL    string         `json:"milestones_url"`
	MirrorURL        string         `json:"mirror_url"`
	NotificationsURL string         `json:"notifications_url"`
	PullsURL         string         `json:"pulls_url"`
	ReleasesURL      string         `json:"releases_url"`
	SSHURL           string         `json:"ssh_url"`
	StargazersURL    string         `json:"stargazers_url"`
	StatusesURL      string         `json:"statuses_url"`
	SubscribersURL   string         `json:"subscribers_url"`
	SubscriptionURL  string         `json:"subscription_url"`
	SvnURL           string         `json:"svn_url"`
	TagsURL          string         `json:"tags_url"`
	TeamsURL         string         `json:"teams_url"`
	TreesURL         string         `json:"trees_url"`
	Homepage         string         `json:"homepage"`
	Language         interface{}    `json:"language"`
	ForksCount       int            `json:"forks_count"`
	StargazersCount  int            `json:"stargazers_count"`
	WatchersCount    int            `json:"watchers_count"`
	Size             int            `json:"size"`
	DefaultBranch    string         `json:"default_branch"`
	OpenIssuesCount  int            `json:"open_issues_count"`
	HasIssues        bool           `json:"has_issues"`
	HasWiki          bool           `json:"has_wiki"`
	HasPages         bool           `json:"has_pages"`
	HasDownloads     bool           `json:"has_downloads"`
	PushedAt         time.Time      `json:"pushed_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	Permissions      GHPermissions  `json:"permissions"`
	SubscribersCount int            `json:"subscribers_count"`
	Organization     GHOrganization `json:"organization"`
}

// GHPermissions represents Github permissions
type GHPermissions struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

// GHOrganization represents a single Github organization
type GHOrganization struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}
