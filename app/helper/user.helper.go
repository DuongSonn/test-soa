package helper

import (
	"context"

	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/repository"
	"sondth-test_soa/utils"

	"github.com/google/uuid"
)

type userHelper struct {
	postgresRepo repository.RepositoryCollections
}

func NewUserHelper(postgresRepo repository.RepositoryCollections) IUserHelper {
	return &userHelper{
		postgresRepo: postgresRepo,
	}
}

func (s *userHelper) IsValidRole(role string) bool {
	return role == entity.ROLE_ADMIN || role == entity.ROLE_USER
}

func (s *userHelper) IsValidPermission(ctx context.Context, userID uuid.UUID) bool {
	user, ok := ctx.Value(string(utils.USER_CONTEXT_KEY)).(*entity.User)
	if !ok {
		return false
	}
	return user.ID == userID && user.IsAdmin()
}
