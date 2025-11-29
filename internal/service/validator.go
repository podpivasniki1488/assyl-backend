package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type userValidator struct {
	repo *repository.Repository
}

func NewUserValidator(repo *repository.Repository) UserValidator {
	return &userValidator{
		repo: repo,
	}
}

func (h *userValidator) Login(ctx context.Context, email, psw string) (*model.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	user, err := h.repo.UserRepo.FindByUsername(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(psw)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return nil, fmt.Errorf("password is wrong")
		default:
			return nil, err
		}
	}

	return user, nil
}
