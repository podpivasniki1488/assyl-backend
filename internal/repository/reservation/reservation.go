package reservation

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type reservationRepository struct {
	db     *gorm.DB
	tracer trace.Tracer
	debug  bool
}

func NewReservationRepository(db *gorm.DB, tracer trace.Tracer) ReservationRepo {
	return &reservationRepository{
		db:     db,
		tracer: tracer,
		debug:  false,
	}
}

func (r *reservationRepository) CreateReservation(ctx context.Context, reservation *model.CinemaReservation) error {
	ctx, span := r.tracer.Start(ctx, "reservationRepository.CreateReservation")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(reservation).Error; err != nil {
		return model.ErrDBUnexpected.WithErr(err)
	}

	return nil
}

func (r *reservationRepository) GetByFilters(ctx context.Context, req *model.CinemaReservation) ([]model.CinemaReservation, error) {
	ctx, span := r.tracer.Start(ctx, "reservationRepository.GetByFilters")
	defer span.End()

	query := r.db.WithContext(ctx)

	if !req.From.IsZero() && !req.To.IsZero() {
		query = query.Where("from < ? AND to > ?", req.From, req.To)
	}

	if !req.From.IsZero() {
		query = query.Where("from >= ?", req.From)
	}

	if !req.To.IsZero() {
		query = query.Where("to <= ?", req.To)
	}

	if req.ID != uuid.Nil {
		query = query.Where("id != ?", req.ID)
	}

	if req.UserID != uuid.Nil {
		query = query.Where("user_id = ?", req.UserID)
	}

	var resp []model.CinemaReservation
	if err := query.Find(&resp).Error; err != nil {
		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	return resp, nil
}

func (r *reservationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.CinemaReservation, error) {
	ctx, span := r.tracer.Start(ctx, "reservationRepository.GetByUserID")
	defer span.End()

	var resp []model.CinemaReservation
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&resp)

	if query.Error != nil {
		return nil, model.ErrDBUnexpected.WithErr(query.Error)
	}

	return resp, nil
}

func (r *reservationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.CinemaReservation, error) {
	ctx, span := r.tracer.Start(ctx, "reservationRepository.GetByID")
	defer span.End()

	var resp model.CinemaReservation
	query := r.db.WithContext(ctx).Where("id = ?", id).First(&resp)

	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return nil, model.ErrReservationNotFound
		}

		return nil, model.ErrDBUnexpected.WithErr(query.Error)
	}

	return &resp, nil
}
