package order

import (
	"context"

	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

type OrderRepo interface {
	Create(ctx context.Context, req *model.Order) error
	GetByFilters(ctx context.Context, req *model.GetOrderRequest) ([]model.Order, error)
}
