package service

import (
	"context"

	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type apartment struct {
	repo  *repository.Repository
	trace trace.Tracer
}

func NewApartmentService(repo *repository.Repository) Apartment {
	return &apartment{repo, otel.Tracer("apartmentService")}
}

func (a *apartment) CreateApartment(ctx context.Context, req model.Apartment) error {
	ctx, span := a.trace.Start(ctx, "apartmentService.CreateApartment")
	defer span.End()

	if _, err := a.repo.ApartmentRepo.CreateApartment(ctx, req); err != nil {
		return err
	}

	return nil
}

func (a *apartment) GetApartment(ctx context.Context, req model.Apartment) (*model.Apartment, error) {
	ctx, span := a.trace.Start(ctx, "apartmentService.GetApartment")
	defer span.End()

	res, err := a.repo.ApartmentRepo.GetApartmentByFloorAndNum(ctx, req.Floor, req.DoorNumber)
	if err != nil {
		return nil, err
	}

	return res, nil
}
