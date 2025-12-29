package reservation

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type reservationRepository struct {
	db     *gorm.DB
	tracer trace.Tracer
	debug  bool
}

func NewReservationRepository(db *gorm.DB) ReservationRepo {
	return &reservationRepository{
		db:     db,
		tracer: otel.Tracer("reservationRepository"),
		debug:  true,
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

func (r *reservationRepository) GetByFilters(ctx context.Context, req *model.GetReservationRequest) ([]model.CinemaReservation, error) {
	ctx, span := r.tracer.Start(ctx, "reservationRepository.GetByFilters")
	defer span.End()

	query := r.db.WithContext(ctx)

	if !req.StartTimeTo.IsZero() {
		query = query.Where("start_time <= ?", req.StartTimeTo)
	}

	if !req.StartTimeFrom.IsZero() {
		query = query.Where("start_time >= ?", req.StartTimeFrom)
	}

	if !req.EndTimeFrom.IsZero() {
		query = query.Where("end_time >= ?", req.EndTimeFrom)
	}

	if !req.EndTimeTo.IsZero() {
		query = query.Where("end_time <= ?", req.EndTimeTo)
	}

	if req.UserID != uuid.Nil {
		query = query.Where("user_id = ?", req.UserID)
	}

	if req.PeopleNumFrom != nil {
		query = query.Where("people_num >= ?", req.PeopleNumFrom)
	}

	if req.PeopleNumTo != nil {
		query = query.Where("people_num <= ?", req.PeopleNumTo)
	}

	if req.IsApproved != nil {
		query = query.Where("is_approved = ?", req.IsApproved)
	}

	if r.debug {
		query = query.Debug()
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

func (r *reservationRepository) ApproveReservation(ctx context.Context, reservationID uuid.UUID) error {
	ctx, span := r.tracer.Start(ctx, "reservationRepository.ApproveReservation")
	defer span.End()

	query := r.db.WithContext(ctx).
		Model(&model.CinemaReservation{}).
		Where("id = ?", reservationID).
		Update("is_approved", true)

	if query.Error != nil {
		return model.ErrDBUnexpected.WithErr(query.Error)
	}

	return nil
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
