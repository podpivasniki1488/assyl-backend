package apartment

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type apartment struct {
	db     *gorm.DB
	tracer trace.Tracer
}

func NewApartmentRepo(db *gorm.DB) ApartmentRepo {
	return &apartment{
		db:     db,
		tracer: otel.Tracer("apartmentRepo"),
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

		return &req, nil

	}

	return &res, nil
}

func (a *apartment) GetApartmentByFloorAndNum(ctx context.Context, floor uint8, doorNum uint16) (*model.Apartment, error) {
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

func (a *apartment) GetApartmentByID(ctx context.Context, id uuid.UUID) (*model.Apartment, error) {
	ctx, span := a.tracer.Start(ctx, "apartmentRepo.GetApartmentByID")
	defer span.End()

	var res model.Apartment
	query := a.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&res)

	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return nil, model.ErrApartmentNotFound
		}

		return nil, model.ErrDBUnexpected.WithErr(query.Error)
	}

	return &res, nil
}

func (a *apartment) UpdateApartment(ctx context.Context, updatedApp *model.Apartment) error {
	ctx, span := a.tracer.Start(ctx, "apartmentRepo.UpdateApartment")
	defer span.End()

	findApp := a.db.
		WithContext(ctx).
		Where("id = ?", updatedApp.Id).
		Updates(updatedApp)

	if findApp.Error != nil {
		return model.ErrDBUnexpected.WithErr(findApp.Error)
	}

	return nil
}
