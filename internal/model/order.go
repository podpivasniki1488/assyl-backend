package model

import (
	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/protopb"
)

type Order struct {
	ID        uuid.UUID         `gorm:"type:uuid;not null;default:gen_random_uuid()" json:"id"`
	OrderType protopb.OrderType `gorm:"type:smallint;not null" json:"order_type"`
	UserID    uuid.UUID         `gorm:"type:uuid;not null" json:"user_id"`
	Text      string            `gorm:"type:varchar;not null" json:"text"`
}

func (Order) TableName() string {
	return "orders"
}

type GetOrderRequest struct {
	ID        uuid.UUID
	OrderType string
	UserID    uuid.UUID
	Text      string
}
