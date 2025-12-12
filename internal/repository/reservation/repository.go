package reservation

import (
	"context"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

type ReservationRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.CinemaReservation, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.CinemaReservation, error)
	CreateReservation(ctx context.Context, reservation *model.CinemaReservation) error
	GetByFilters(ctx context.Context, req *model.CinemaReservation) ([]model.CinemaReservation, error)
}
