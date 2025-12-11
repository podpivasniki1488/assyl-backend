package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	UserValidator  UserValidator
	UserManagement UserManagement
	Auth           Auth
	Apartment      Apartment
}

type Auth interface {
	Login(ctx context.Context, user model.User) (token string, err error)
	Register(ctx context.Context, user model.User) error
	Confirm(ctx context.Context, username, otpCode string) error
}

type UserValidator interface {
	Login(ctx context.Context, login, psw string) (*model.User, error)
}

type UserManagement interface {
	DeleteUserByEmail(ctx context.Context, email string) error
	BindApartmentToUser(ctx context.Context, username string, apartmentId uuid.UUID) error
}

type Apartment interface {
	CreateApartment(ctx context.Context, req model.Apartment) error
	GetApartment(ctx context.Context, req model.Apartment) (*model.Apartment, error)
}

func NewService(
	repo *repository.Repository,
	tracer trace.Tracer,
	redisCli *redis.Client,
	secretKey string,
) *Service {
	return &Service{
		UserValidator:  NewUserValidator(repo),
		UserManagement: NewUserManagement(repo, tracer),
		Auth:           NewAuthService(repo, secretKey, tracer, redisCli),
		Apartment:      NewApartmentService(repo, tracer),
	}
}
