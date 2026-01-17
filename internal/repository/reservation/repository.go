package reservation

import (
	"context"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

type ReservationRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.CinemaReservation, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.CinemaReservation, error)
	CreateReservationsWithSlots(ctx context.Context, res *model.CinemaReservation, dailySlotIDs []uint64) error
	GetByFilters(ctx context.Context, req *model.GetReservationRequest) ([]model.CinemaReservation, error)
	ApproveReservation(ctx context.Context, reservationID uuid.UUID) error
}
