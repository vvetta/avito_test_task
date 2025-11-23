package prrepository

import (
	"context"
	"errors"

	"avito_test_task/internal/domain"

	"gorm.io/gorm"
)

type PrRepository struct {
	db *gorm.DB
}

func NewPrRepository(db *gorm.DB) *PrRepository {
	return &PrRepository{
		db: db,
	}
}

func (p *PrRepository) Create(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error) {
	prModel := domainToModel(pr)
	reviewerModels := reviewersToModels(pr.PullRequestID, pr.AssignedReviewers)

	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&prModel).Error; err != nil {
			return err
		}
		if len(reviewerModels) > 0 {
			if err := tx.Create(&reviewerModels).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return domain.PullRequest{}, err
	}

	return p.GetByID(ctx, pr.PullRequestID)
}

func (p *PrRepository) GetByID(ctx context.Context, prID string) (domain.PullRequest, error) {
	var prModel PullRequestModel
	if err := p.db.WithContext(ctx).
		Where("pull_request_id = ?", prID).
		First(&prModel).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.PullRequest{}, domain.ErrPullRequestNotFound
		}
		return domain.PullRequest{}, err
	}

	var reviewers []PullRequestReviewerModel
	if err := p.db.WithContext(ctx).
		Where("pull_request_id = ?", prID).
		Find(&reviewers).Error; err != nil {
		return domain.PullRequest{}, err
	}

	return modelToDomain(prModel, reviewers), nil
}

func (p *PrRepository) Update(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error) {
	prModel := domainToModel(pr)
	reviewerModels := reviewersToModels(pr.PullRequestID, pr.AssignedReviewers)

	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&prModel).Error; err != nil {
			return err
		}

		if err := tx.
			Where("pull_request_id = ?", pr.PullRequestID).
			Delete(&PullRequestReviewerModel{}).Error; err != nil {
			return err
		}

		if len(reviewerModels) > 0 {
			if err := tx.Create(&reviewerModels).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return domain.PullRequest{}, err
	}

	return p.GetByID(ctx, pr.PullRequestID)
}

func (p *PrRepository) GetByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequest, error) {
	var links []PullRequestReviewerModel
	if err := p.db.WithContext(ctx).
		Where("reviewer_id = ?", reviewerID).
		Find(&links).Error; err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return []domain.PullRequest{}, nil
	}

	prIDs := make([]string, 0, len(links))
	for _, l := range links {
		prIDs = append(prIDs, l.PullRequestID)
	}

	var prModels []PullRequestModel
	if err := p.db.WithContext(ctx).
		Where("pull_request_id IN ?", prIDs).
		Find(&prModels).Error; err != nil {
		return nil, err
	}

	result := make([]domain.PullRequest, 0, len(prModels))
	for _, m := range prModels {
		var reviewers []PullRequestReviewerModel
		if err := p.db.WithContext(ctx).
			Where("pull_request_id = ?", m.PullRequestID).
			Find(&reviewers).Error; err != nil {
			return nil, err
		}
		result = append(result, modelToDomain(m, reviewers))
	}

	return result, nil
}

