package prrepository

import (
	"time"
)

type PullRequestModel struct {
    PullRequestID   string `gorm:"primaryKey"`
    PullRequestName string
    AuthorID        string `gorm:"index"`
    Status          string
    CreatedAt       *time.Time
    MergedAt        *time.Time
}

func (PullRequestModel) TableName() string { return "pull_requests" }

type PullRequestReviewerModel struct {
    PullRequestID string `gorm:"primaryKey"`
    ReviewerID    string `gorm:"primaryKey"`
}

func (PullRequestReviewerModel) TableName() string { return "pull_request_reviewers" }
