package apartment

import (
	"context"

	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

type ApartmentRepo interface {
	CreateApartment(ctx context.Context, req model.Apartment) (*model.Apartment, error)
	GetApartment(ctx context.Context, floor uint8, doorNum uint16) (*model.Apartment, error)
}
