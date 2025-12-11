package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository"
	"github.com/podpivasniki1488/assyl-backend/protopb"
	"go.opentelemetry.io/otel/trace"
)

type reservation struct {
	tracer trace.Tracer
	repo   *repository.Repository
}

func NewReservation(repo *repository.Repository, tracer trace.Tracer) Reservation {
	return &reservation{
		tracer: tracer,
		repo:   repo,
	}
}

func (r *reservation) GetFilteredReservations(ctx context.Context, req model.CinemaReservation, userID uuid.UUID) ([]model.CinemaReservation, error) {
	ctx, span := r.tracer.Start(ctx, "reservation.GetReservation")
	defer span.End()

	user, err := r.repo.UserRepo.FindById(ctx, userID)
	if err != nil {
		return nil, err
	}

	reservations, err := r.repo.ReservationRepo.GetByFilters(ctx, &model.CinemaReservation{
		From: req.From,
		To:   req.To,
	})
	if err != nil {
		return nil, err
	}

	return r.filterReservation(ctx, reservations, *user)
}

func (r *reservation) filterReservation(ctx context.Context, reservations []model.CinemaReservation, user model.User) ([]model.CinemaReservation, error) {
	ctx, span := r.tracer.Start(ctx, "reservation.filterReservation")
	defer span.End()

	if user.RoleID == protopb.Role_ADMIN || user.RoleID == protopb.Role_GOD {
		return reservations, nil
	}

	res := make([]model.CinemaReservation, 0, len(reservations))
	for _, resv := range reservations {
		if resv.UserID == user.ID {
			res = append(res, resv)
		}
	}

	return res, nil
}

func (r *reservation) MakeReservation(ctx context.Context, req *model.CinemaReservation) error {
	ctx, span := r.tracer.Start(ctx, "reservation.MakeReservation")
	defer span.End()

	reservations, err := r.repo.ReservationRepo.GetByFilters(ctx, &model.CinemaReservation{
		From: req.From,
		To:   req.To,
	})
	if err != nil {
		return err
	}

	if len(reservations) != 0 {
		return model.ErrCinemaBusy
	}

	if err = r.repo.ReservationRepo.CreateReservation(ctx, req); err != nil {
		return err
	}

	return nil
}
