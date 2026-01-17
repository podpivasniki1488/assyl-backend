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
	Feedback       Feedback
	Order          Order
}

type Auth interface {
	Login(ctx context.Context, user model.User) (token string, u *model.User, err error)
	Register(ctx context.Context, user model.User) error
	Confirm(ctx context.Context, username, otpCode string) error
}

type UserValidator interface {
	Login(ctx context.Context, login, psw string) (*model.User, error)
}

type UserManagement interface {
	DeleteUserByUsername(ctx context.Context, username string) error
	BindApartmentToUser(ctx context.Context, username string, apartmentId uuid.UUID) error
}

type Apartment interface {
	CreateApartment(ctx context.Context, req model.Apartment) error
	GetApartment(ctx context.Context, req model.Apartment) (*model.Apartment, error)
}

type Reservation interface {
	MakeReservation(ctx context.Context, userID uuid.UUID, date time.Time, positions []int16, peopleNum uint8, role, username string) error
	GetUserReservations(ctx context.Context, req model.CinemaReservation, start, end time.Time) ([]model.CinemaReservation, error)
	GetUnfilteredReservations(ctx context.Context, req model.CinemaReservation) ([]model.CinemaReservation, error)
	ApproveReservation(ctx context.Context, id uuid.UUID) error
	GetFreeSlots(ctx context.Context, date time.Time) ([]model.DailySlot, [][2]int16, error)
}

type Channel interface {
	SendChannelMessage(ctx context.Context, msg model.ChannelMessage) error
	GetByTimePeriod(ctx context.Context, from, to time.Time) ([]model.ChannelMessage, error)
}

type Feedback interface {
	CreateFeedback(ctx context.Context, req model.Feedback) error
	GetFeedbacks(ctx context.Context, req model.GetFeedbackRequest) ([]model.Feedback, error)
}

type Order interface {
	OrderService(ctx context.Context, req *model.Order) error
	GetUserOrders(ctx context.Context, req *model.GetOrderRequest, role string) ([]model.Order, error)
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
		Feedback:       NewFeedback(repo),
		Order:          NewOrderService(repo),
	}
}
