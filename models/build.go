package models

import "time"

// UnparsedBuild represents the unparsed response from the CI
type UnparsedBuild struct {
	ID           int    `json:"id"`
	Number       int    `json:"number"`
	Event        string `json:"event"`
	Status       string `json:"status"`
	EnqueuedAt   int64  `json:"enqueued_at"`
	CreatedAt    int64  `json:"created_at"`
	StartedAt    int64  `json:"started_at"`
	FinishedAt   int64  `json:"finished_at"`
	DeployTo     string `json:"deploy_to"`
	Commit       string `json:"commit"`
	Branch       string `json:"branch"`
	Ref          string `json:"ref"`
	Refspec      string `json:"refspec"`
	Remote       string `json:"remote"`
	Title        string `json:"title"`
	Message      string `json:"message"`
	Timestamp    int64  `json:"timestamp"`
	Author       string `json:"author"`
	AuthorAvatar string `json:"author_avatar"`
	AuthorEmail  string `json:"author_email"`
	LinkURL      string `json:"link_url"`
}

// UnparsedBuilds represents a list of UnparsedBuild
type UnparsedBuilds []UnparsedBuild

// Build is the parsed CI response
type Build struct {
	ID           int
	Number       int
	Event        string
	Status       string
	EnqueuedAt   time.Time
	CreatedAt    time.Time
	StartedAt    time.Time
	FinishedAt   time.Time
	DeployTo     string
	Commit       string
	Branch       string
	Ref          string
	Refspec      string
	Remote       string
	Title        string
	Message      string
	Timestamp    time.Time
	Author       string
	AuthorAvatar string
	AuthorEmail  string
	LinkURL      string
	Duration     time.Duration
}

// Parse parses the times in the UnparsedBuild and returns a Build
func (u UnparsedBuild) Parse() Build {
	b := Build{
		ID:           u.ID,
		Number:       u.Number,
		Event:        u.Event,
		Status:       u.Status,
		EnqueuedAt:   time.Unix(u.EnqueuedAt, 0),
		CreatedAt:    time.Unix(u.CreatedAt, 0),
		StartedAt:    time.Unix(u.StartedAt, 0),
		FinishedAt:   time.Unix(u.FinishedAt, 0),
		DeployTo:     u.DeployTo,
		Commit:       u.Commit,
		Branch:       u.Branch,
		Ref:          u.Ref,
		Refspec:      u.Refspec,
		Remote:       u.Remote,
		Title:        u.Title,
		Message:      u.Message,
		Timestamp:    time.Unix(u.Timestamp, 0),
		Author:       u.Author,
		AuthorAvatar: u.AuthorAvatar,
		AuthorEmail:  u.AuthorEmail,
		LinkURL:      u.LinkURL,
	}
	b.Duration = b.FinishedAt.Sub(b.StartedAt)
	return b
}

// Builds represent multiple builds
type Builds []Build
