package apartment

import (
	"context"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

type ApartmentRepo interface {
	CreateApartment(ctx context.Context, req model.Apartment) (*model.Apartment, error)
	GetApartmentByFloorAndNum(ctx context.Context, floor uint8, doorNum uint16) (*model.Apartment, error)
	GetApartmentByID(ctx context.Context, id uuid.UUID) (*model.Apartment, error)
	UpdateApartment(ctx context.Context, updatedApp *model.Apartment) error
}
