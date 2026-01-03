package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository"
	"github.com/podpivasniki1488/assyl-backend/protopb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type order struct {
	repo   *repository.Repository
	tracer trace.Tracer
}

func NewOrderService(repo *repository.Repository) Order {
	return &order{
		repo:   repo,
		tracer: otel.Tracer("orderService"),
	}
}

func (o *order) OrderService(ctx context.Context, req *model.Order) error {
	ctx, span := o.tracer.Start(ctx, "OrderService.Order")
	defer span.End()

	if err := o.repo.OrderRepo.Create(ctx, req); err != nil {
		return err
	}

	// TODO: send notification via whatsapp to admin's phone num

	return nil
}

func (o *order) GetUserOrders(ctx context.Context, req *model.GetOrderRequest, role string) ([]model.Order, error) {
	ctx, span := o.tracer.Start(ctx, "GetUserOrders.Order")
	defer span.End()

	if role == protopb.Role_GOD.String() || role == protopb.Role_ADMIN.String() {
		req.UserID = uuid.Nil
	}

	resp, err := o.repo.OrderRepo.GetByFilters(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
