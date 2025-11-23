package domain

import (
	"time"
)

type TeamMember struct {
	UserID string
	Username string
	IsActive bool
}

type Team struct {
	TeamName string
	Members []TeamMember	
}

type User struct {
	UserID string
	TeamName string
	Username string
	IsActive bool
}

type PullRequestStatus string

const (
	StatusOpen   PullRequestStatus = "OPEN"
	StatusMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID string
	PullRequestName string
	AuthorID string
	Status PullRequestStatus
	AssignedReviewers []string
	CreatedAt *time.Time
	MergedAt *time.Time	
}
