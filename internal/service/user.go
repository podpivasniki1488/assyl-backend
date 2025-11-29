package service

import (
	"context"

	"github.com/podpivasniki1488/assyl-backend/internal/repository"
)

type userManagement struct {
	repo *repository.Repository
}

func NewUserManagement(repo *repository.Repository) UserManagement {
	return &userManagement{repo}
}

func (u *userManagement) DeleteUserByEmail(ctx context.Context, email string) error {
	if err := u.repo.UserRepo.DeleteByUsername(ctx, email); err != nil {
		return err
	}

	return nil
}
