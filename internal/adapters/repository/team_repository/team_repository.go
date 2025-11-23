package teamrepository

import (
	"context"
	"errors"

	"avito_test_task/internal/domain"

	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (t *TeamRepository) Create(ctx context.Context, team domain.Team) (domain.Team, error) {
model := TeamModel{
		TeamName: team.TeamName,
	}

	if err := t.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Team{}, err
	}

	return domain.Team{
		TeamName: model.TeamName,
		Members:  team.Members,
	}, nil
}

func (t *TeamRepository) GetByName(ctx context.Context, name string) (domain.Team, error) {
	var model TeamModel
	err := t.db.WithContext(ctx).
		Where("team_name = ?", name).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Team{}, domain.ErrTeamNotFound
		}
		return domain.Team{}, err
	}

	return domain.Team{
		TeamName: model.TeamName,
		Members:  nil,
	}, nil
}

