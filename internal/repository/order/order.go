package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/protopb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type orderRepository struct {
	db     *gorm.DB
	tracer trace.Tracer
	debug  bool
}

func NewOrderRepository(db *gorm.DB) OrderRepo {
	return &orderRepository{
		db:     db,
		tracer: otel.Tracer("orderRepository"),
		debug:  true,
	}
}

func (o *orderRepository) Create(ctx context.Context, req *model.Order) error {
	ctx, span := o.tracer.Start(ctx, "orderRepository.Create")
	defer span.End()

	query := o.db.Create(&req)

	if o.debug {
		query = query.Debug()
	}

	if err := query.Error; err != nil {
		return model.ErrDBUnexpected.WithErr(err)
	}

	return nil
}

func (o *orderRepository) GetByFilters(ctx context.Context, req *model.GetOrderRequest) ([]model.Order, error) {
	ctx, span := o.tracer.Start(ctx, "orderRepository.GetByFilters")
	defer span.End()

	query := o.db.WithContext(ctx)

	if val, ok := protopb.OrderType_value[req.OrderType]; ok {
		query = query.Where("order_type = ?", val)
	}

	if req.UserID != uuid.Nil {
		query = query.Where("user_id = ?", req.UserID)
	}

	if req.ID != uuid.Nil {
		query = query.Where("order_id = ?", req.ID)
	}

	if req.Text != "" {
		query = query.Where(fmt.Sprintf("text ilike '%s'", req.Text))
	}

	if o.debug {
		query = query.Debug()
	}

	var orders []model.Order
	if err := query.Find(&orders).Error; err != nil {
		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	return orders, nil
}
