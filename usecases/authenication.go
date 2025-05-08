package usecases

import (
	"context"
	"errors"

	"github.com/luckydevil2007/go-lessons/entities"
)

type UserStorage interface {
	GetByLogin(ctx context.Context, login string) (*entities.User, error)
}

type AuthUseCase struct {
	userStorage UserStorage
}

func NewAuthUseCase(userStorage UserStorage) *AuthUseCase {
	return &AuthUseCase{userStorage: userStorage}
}

func (uc *AuthUseCase) Authenticate(ctx context.Context, username, password string) (int, error) {
	user, _ := uc.userStorage.GetByLogin(ctx, username)
	if user == nil {
		return -1, errors.New("unauthorized")
	}

	return user.ID, nil
}
