package userrepository

import "avito_test_task/internal/domain"

func modelToDomain(m UserModel) domain.User {
	return domain.User{
		UserID:   m.UserID,
		TeamName: m.TeamName,
		Username: m.Username,
		IsActive: m.IsActive,
	}
}

