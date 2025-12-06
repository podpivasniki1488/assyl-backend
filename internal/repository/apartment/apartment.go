package apartment

import (
	"context"
	"errors"

	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type apartment struct {
	db     *gorm.DB
	tracer trace.Tracer
}

func NewApartmentRepo(db *gorm.DB, tracer trace.Tracer) ApartmentRepo {
	return &apartment{
		db:     db,
		tracer: tracer,
	}
}

func (a *apartment) CreateApartment(ctx context.Context, req model.Apartment) (*model.Apartment, error) {
	ctx, span := a.tracer.Start(ctx, "apartmentRepo.CreateApartment")
	defer span.End()

	var res model.Apartment
	query := a.db.
		WithContext(ctx).
		Where("floor = ?", req.Floor).
		Where("door_number = ?", req.DoorNumber).
		First(&res)

	if err := query.Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrDBUnexpected.WithErr(query.Error)
		}

		if err := a.db.WithContext(ctx).Create(&req).Error; err != nil {
			return nil, model.ErrDBUnexpected.WithErr(err)
		}

	}
	return &res, nil
}

func (a *apartment) GetApartment(ctx context.Context, floor uint8, doorNum uint16) (*model.Apartment, error) {
	ctx, span := a.tracer.Start(ctx, "apartmentRepo.GetApartment")
	defer span.End()

	var res model.Apartment
	query := a.db.
		WithContext(ctx).
		Where("floor = ? AND door_number = ?", floor, doorNum).
		First(&res)

	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return nil, model.ErrApartmentNotFound
		}

		return nil, model.ErrDBUnexpected.WithErr(query.Error)
	}

	return &res, nil
}
