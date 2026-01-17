package model

import (
	"time"

	"github.com/google/uuid"
)

type SlotTemplate struct {
	ID        int16  `gorm:"primaryKey;type:smallserial" json:"id"`
	Code      string `gorm:"type:text;uniqueIndex;not null" json:"code"`
	StartTime string `gorm:"type:time;not null" json:"start_time"` // можно time.Time, но аккуратно с датой; проще string "10:00:00"
	EndTime   string `gorm:"type:time;not null" json:"end_time"`
	Position  int16  `gorm:"type:smallint;uniqueIndex;not null" json:"position"`
	IsActive  bool   `gorm:"type:boolean;not null;default:true" json:"is_active"`
}

func (SlotTemplate) TableName() string { return "slot_templates" }

type DailySlot struct {
	ID         int64     `gorm:"primaryKey;type:bigserial" json:"id"`
	SlotDate   time.Time `gorm:"type:date;not null;index:ix_daily_slots_date" json:"slot_date"`
	TemplateID int16     `gorm:"type:smallint;not null" json:"template_id"`
	Position   int16     `gorm:"type:smallint;not null" json:"position"`
	StartAt    time.Time `gorm:"type:timestamptz;not null" json:"start_at"`
	EndAt      time.Time `gorm:"type:timestamptz;not null" json:"end_at"`
	IsEnabled  bool      `gorm:"type:boolean;not null;default:true" json:"is_enabled"`
}

func (DailySlot) TableName() string { return "daily_slots" }

type ReservationSlot struct {
	ReservationID uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"reservation_id"`
	DailySlotID   uint64    `gorm:"type:bigint;not null;primaryKey" json:"daily_slot_id"`
}

func (ReservationSlot) TableName() string { return "reservation_slots" }
