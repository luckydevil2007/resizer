package controllers

import (
	"context"

	"github.com/luckydevil2007/audionotes/adapters/repositories"
	"github.com/luckydevil2007/audionotes/producers"
	"github.com/luckydevil2007/audionotes/usecases"
)

type NoteController struct {
	checkAuth *usecases.AuthUseCase
	note      *usecases.NoteUseCase
	repo      *repositories.Repository
	producer  *producers.EventProducer
}

type AuthController struct {
	checkAuth *usecases.AuthUseCase
}

func NewAuthController(authUC *usecases.AuthUseCase) *AuthController {
	return &AuthController{
		checkAuth: authUC,
	}
}

func (c *AuthController) Authenticate(ctx context.Context, username string) context.Context {
	userID, err := c.checkAuth.Authenticate(ctx, username)
	if err != nil {
		return nil
	}

	return context.WithValue(ctx, "id", userID)
}
