package models

// Build represents the builds
type Build struct {
	ID           int    `json:"id"`
	Number       int    `json:"number"`
	Event        string `json:"event"`
	Status       string `json:"status"`
	EnqueuedAt   int    `json:"enqueued_at"`
	CreatedAt    int    `json:"created_at"`
	StartedAt    int    `json:"started_at"`
	FinishedAt   int    `json:"finished_at"`
	DeployTo     string `json:"deploy_to"`
	Commit       string `json:"commit"`
	Branch       string `json:"branch"`
	Ref          string `json:"ref"`
	Refspec      string `json:"refspec"`
	Remote       string `json:"remote"`
	Title        string `json:"title"`
	Message      string `json:"message"`
	Timestamp    int    `json:"timestamp"`
	Author       string `json:"author"`
	AuthorAvatar string `json:"author_avatar"`
	AuthorEmail  string `json:"author_email"`
	LinkURL      string `json:"link_url"`
}

// Builds represent multiple builds
type Builds []Build
