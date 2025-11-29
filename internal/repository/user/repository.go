package user

import (
	"context"

	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

type UserRepo interface {
	FindById(ctx context.Context, id int64) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	DeleteByUsername(ctx context.Context, username string) error
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
}
