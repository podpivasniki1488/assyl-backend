package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	UserValidator  UserValidator
	UserManagement UserManagement
	Auth           Auth
	Apartment      Apartment
	Reservation    Reservation
	Channel        Channel
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

type Reservation interface {
	MakeReservation(ctx context.Context, req *model.CinemaReservation, role, username string) error
	GetUserReservations(ctx context.Context, req model.CinemaReservation) ([]model.CinemaReservation, error)
	GetUnfilteredReservations(ctx context.Context, req model.CinemaReservation) ([]model.CinemaReservation, error)
	ApproveReservation(ctx context.Context, id uuid.UUID) error
}

type Channel interface {
	SendChannelMessage(ctx context.Context, msg model.ChannelMessage) error
	GetByTimePeriod(ctx context.Context, from, to time.Time) ([]model.ChannelMessage, error)
}

func NewService(
	repo *repository.Repository,
	redisCli *redis.Client,
	secretKey string,
) *Service {
	return &Service{
		UserValidator:  NewUserValidator(repo),
		UserManagement: NewUserManagement(repo),
		Auth:           NewAuthService(repo, secretKey, redisCli),
		Apartment:      NewApartmentService(repo),
		Reservation:    NewReservation(repo),
		Channel:        NewChannelService(repo),
	}
}
