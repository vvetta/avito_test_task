package usecase

import (
	"context"

	"avito_test_task/internal/domain"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team domain.Team) (domain.Team, error)
	GetTeam(ctx context.Context, name string) (domain.Team, error)
}

type TeamRepository interface {
	Create(ctx context.Context, team domain.Team) (domain.Team, error)
	GetByName(ctx context.Context, name string) (domain.Team, error)
}

type UserService interface {
	SetUserActive(ctx context.Context, userID string, isActive bool) (domain.User, error)
	GetUserReviews(ctx context.Context, userID string) ([]domain.PullRequest, error) 
}

type UserRepository interface {
	SetActive(ctx context.Context, userID string, isActive bool) (domain.User, error)
	GetByID(ctx context.Context, userID string) (domain.User, error)
	GetTeamMembers(ctx context.Context, teamName string) ([]domain.User, error)
	CreateOrUpdateMany(ctx context.Context, users []domain.User) error
}

type PullRequestService interface {
	CreatePR(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error )
	MergePR(ctx context.Context, prID string) (domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (domain.PullRequest, string, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error)
	GetByID(ctx context.Context, prID string) (domain.PullRequest, error)
	Update(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error)
	GetByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequest, error)
}
