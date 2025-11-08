package models

import "time"

type Site struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Slug         string    `json:"slug"`
	GithubRepo   string    `json:"github_repo"`
	GithubBranch string    `json:"github_branch"`
	Subdirectory string    `json:"subdirectory"`
	CreatedAt    time.Time `json:"created_at"`
}
