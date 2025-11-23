package userrepository

import (
	"context"
	"errors"

	"avito_test_task/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) SetActive(ctx context.Context, userID string, isActive bool) (domain.User, error) {
	var model UserModel

	err := u.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	model.IsActive = isActive

	if err := u.db.WithContext(ctx).Save(&model).Error; err != nil {
		return domain.User{}, err
	}

	return modelToDomain(model), nil
}

func (u *UserRepository) GetByID(ctx context.Context, userID string) (domain.User, error) {
	var model UserModel

	err := u.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return modelToDomain(model), nil
}

func (u *UserRepository) GetTeamMembers(ctx context.Context, teamName string) ([]domain.User, error) {
	var models []UserModel

	if err := u.db.WithContext(ctx).
		Where("team_name = ?", teamName).
		Find(&models).Error; err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(models))
	for _, m := range models {
		users = append(users, modelToDomain(m))
	}

	return users, nil
}

func (u *UserRepository) CreateOrUpdateMany(ctx context.Context, users []domain.User) error {
	if len(users) == 0 {
		return nil
	}

	for _, du := range users {
		model := UserModel{
			UserID:   du.UserID,
			TeamName: du.TeamName,
			Username: du.Username,
			IsActive: du.IsActive,
		}

		if err := u.db.WithContext(ctx).Save(&model).Error; err != nil {
			return err
		}
	}

	return nil
}

