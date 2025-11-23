package usecase

import (
	"context"
	"errors"

	"avito_test_task/internal/domain"
)

type teamService struct {
	teamRepo TeamRepository
	userRepo UserRepository
}

func NewTeamService(teamRepo TeamRepository, userRepo UserRepository) TeamService {
	return &teamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *teamService) CreateTeam(ctx context.Context, team domain.Team) (domain.Team, error) {
	existing, err := s.teamRepo.GetByName(ctx, team.TeamName)
	if err == nil && existing.TeamName != "" {
		return domain.Team{}, domain.ErrTeamExists
	}
	if err != nil && !errors.Is(err, domain.ErrTeamNotFound) {
		return domain.Team{}, err
	}

	users := make([]domain.User, 0, len(team.Members))
	for _, m := range team.Members {
		users = append(users, domain.User{
			UserID:   m.UserID,
			Username: m.Username,
			TeamName: team.TeamName,
			IsActive: m.IsActive,
		})
	}

	if err := s.userRepo.CreateOrUpdateMany(ctx, users); err != nil {
		return domain.Team{}, err
	}

	created, err := s.teamRepo.Create(ctx, domain.Team{
		TeamName: team.TeamName,
	})
	if err != nil {
		return domain.Team{}, err
	}

	members := make([]domain.TeamMember, 0, len(users))
	for _, u := range users {
		members = append(members, domain.TeamMember{
			UserID:   u.UserID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}
	created.Members = members

	return created, nil
}

func (s *teamService) GetTeam(ctx context.Context, teamName string) (domain.Team, error) {
	team, err := s.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		return domain.Team{}, err
	}

	users, err := s.userRepo.GetTeamMembers(ctx, teamName)
	if err != nil {
		return domain.Team{}, err
	}

	members := make([]domain.TeamMember, 0, len(users))
	for _, u := range users {
		members = append(members, domain.TeamMember{
			UserID:   u.UserID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}
	team.Members = members

	return team, nil
}
