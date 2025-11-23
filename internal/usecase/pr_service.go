package usecase

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"avito_test_task/internal/domain"
)

type prService struct {
	prRepo   PullRequestRepository
	userRepo UserRepository
	teamRepo TeamRepository
}

func NewPRService(
	prRepo PullRequestRepository,
	userRepo UserRepository,
	teamRepo TeamRepository,
) PullRequestService {
	return &prService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *prService) CreatePR(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error) {
	existing, err := s.prRepo.GetByID(ctx, pr.PullRequestID)
	if err == nil && existing.PullRequestID != "" {
		return domain.PullRequest{}, domain.ErrPullRequestExists
	}
	if err != nil && !errors.Is(err, domain.ErrPullRequestNotFound) {
		return domain.PullRequest{}, err
	}

	author, err := s.userRepo.GetByID(ctx, pr.AuthorID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	members, err := s.userRepo.GetTeamMembers(ctx, author.TeamName)
	if err != nil {
		return domain.PullRequest{}, err
	}

	candidates := make([]string, 0, len(members))
	for _, m := range members {
		if !m.IsActive {
			continue
		}
		if m.UserID == pr.AuthorID {
			continue
		}
		candidates = append(candidates, m.UserID)
	}

	assigned := pickReviewers(candidates, 2)

	now := time.Now().UTC()
	newPR := domain.PullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            domain.StatusOpen,
		AssignedReviewers: assigned,
		CreatedAt:         &now,
		MergedAt:          nil,
	}

	return s.prRepo.Create(ctx, newPR)
}

func (s *prService) MergePR(ctx context.Context, prID string) (domain.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	if pr.Status == domain.StatusMerged {
		if pr.MergedAt == nil {
			now := time.Now().UTC()
			pr.MergedAt = &now
			pr, err = s.prRepo.Update(ctx, pr)
			if err != nil {
				return domain.PullRequest{}, err
			}
		}
		return pr, nil
	}

	now := time.Now().UTC()
	pr.Status = domain.StatusMerged
	pr.MergedAt = &now

	return s.prRepo.Update(ctx, pr)
}

func (s *prService) ReassignReviewer(ctx context.Context, prID, oldUserID string) (domain.PullRequest, string, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	if pr.Status == domain.StatusMerged {
		return domain.PullRequest{}, "", domain.ErrPullRequestMerged
	}

	idx := -1
	for i, id := range pr.AssignedReviewers {
		if id == oldUserID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return domain.PullRequest{}, "", domain.ErrReviewerNotAssigned
	}

	reviewer, err := s.userRepo.GetByID(ctx, oldUserID)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	members, err := s.userRepo.GetTeamMembers(ctx, reviewer.TeamName)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	candidates := make([]string, 0, len(members))
	for _, m := range members {
		if !m.IsActive || m.UserID == oldUserID {
			continue
		}
		already := false
		for _, assigned := range pr.AssignedReviewers {
			if assigned == m.UserID {
				already = true
				break
			}
		}
		if already {
			continue
		}
		candidates = append(candidates, m.UserID)
	}

	if len(candidates) == 0 {
		return domain.PullRequest{}, "", domain.ErrNoCandidate
	}

	newIdx := rand.Intn(len(candidates))
	newID := candidates[newIdx]

	pr.AssignedReviewers[idx] = newID

	updated, err := s.prRepo.Update(ctx, pr)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	return updated, newID, nil
}

func pickReviewers(candidates []string, max int) []string {
	if len(candidates) == 0 {
		return nil
	}
	if len(candidates) <= max {
		perm := rand.Perm(len(candidates))
		res := make([]string, 0, len(candidates))
		for _, idx := range perm {
			res = append(res, candidates[idx])
		}
		return res
	}

	perm := rand.Perm(len(candidates))[:max]
	res := make([]string, 0, max)
	for _, idx := range perm {
		res = append(res, candidates[idx])
	}
	return res
}
