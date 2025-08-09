package usecases

import (
	"context"
	"errors"

	"github.com/luckydevil2007/audionotes/entities"
)

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (*entities.User, error)
}

type AuthUseCase struct {
	userRepo UserRepository
}

func NewAuthUseCase(userRepo UserRepository) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo}
}

func (uc *AuthUseCase) Authenticate(ctx context.Context, username string) (int, error) {
	user, _ := uc.userRepo.GetByLogin(ctx, username)
	if user == nil {
		return -1, errors.New("unauthorized")
	}

	return user.ID, nil
}
