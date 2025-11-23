package http

import (
	"avito_test_task/internal/adapters/http/openapi"
	"avito_test_task/internal/domain"
)

func fromAPITeam(openapiTeam openapi.Team) domain.Team {
	domainMembers := make([]domain.TeamMember, len(openapiTeam.Members))

	for i, m := range openapiTeam.Members {
		domainMembers[i] = domain.TeamMember{
			UserID:   m.UserId,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	return domain.Team{
		TeamName: openapiTeam.TeamName,
		Members:  domainMembers,
	}
}

func toAPITeam(domainTeam domain.Team) openapi.Team {
	apiMembers := make([]openapi.TeamMember, len(domainTeam.Members))

	for i, m := range domainTeam.Members {
		apiMembers[i] = openapi.TeamMember{
			UserId:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	return openapi.Team{
		TeamName: domainTeam.TeamName,
		Members:  apiMembers,
	}
}

func toAPIUser(u domain.User) openapi.User {
	return openapi.User{
		UserId:   u.UserID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func toAPIPullRequest(pr domain.PullRequest) openapi.PullRequest {
	return openapi.PullRequest{
		PullRequestId:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorId:          pr.AuthorID,
		Status:            openapi.PullRequestStatus(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func toAPIPullRequestShorts(prs []domain.PullRequest) []openapi.PullRequestShort {
	result := make([]openapi.PullRequestShort, 0, len(prs))
	for _, pr := range prs {
		result = append(result, openapi.PullRequestShort{
			PullRequestId:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorId:        pr.AuthorID,
			Status:          openapi.PullRequestShortStatus(pr.Status),
		})
	}
	return result
}
