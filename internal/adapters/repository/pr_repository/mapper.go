package prrepository

import "avito_test_task/internal/domain"

func modelToDomain(pr PullRequestModel, reviewers []PullRequestReviewerModel) domain.PullRequest {
	reviewerIDs := make([]string, 0, len(reviewers))
	for _, r := range reviewers {
		reviewerIDs = append(reviewerIDs, r.ReviewerID)
	}

	return domain.PullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            domain.PullRequestStatus(pr.Status),
		AssignedReviewers: reviewerIDs,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func domainToModel(pr domain.PullRequest) PullRequestModel {
	return PullRequestModel{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          string(pr.Status),
		CreatedAt:       pr.CreatedAt,
		MergedAt:        pr.MergedAt,
	}
}

func reviewersToModels(prID string, reviewers []string) []PullRequestReviewerModel {
	result := make([]PullRequestReviewerModel, 0, len(reviewers))
	for _, id := range reviewers {
		result = append(result, PullRequestReviewerModel{
			PullRequestID: prID,
			ReviewerID:    id,
		})
	}
	return result
}

