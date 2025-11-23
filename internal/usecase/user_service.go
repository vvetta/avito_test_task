package usecase

import (
	"context"

	"avito_test_task/internal/domain"
)

type userService struct {
	userRepo UserRepository
	prRepo   PullRequestRepository
}

func NewUserService(userRepo UserRepository, prRepo PullRequestRepository) UserService {
	return &userService{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

func (s *userService) SetUserActive(ctx context.Context, userID string, isActive bool) (domain.User, error) {
	return s.userRepo.SetActive(ctx, userID, isActive)
}

func (s *userService) GetUserReviews(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	prs, err := s.prRepo.GetByReviewer(ctx, userID)
	if err != nil {
		return nil, err
	}
	return prs, nil
}
