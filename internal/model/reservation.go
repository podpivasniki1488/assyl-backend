package model

import (
	"time"

	"github.com/google/uuid"
)

type CinemaReservation struct {
	ID         uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;foreignKey:UserID;references:ID" json:"user_id"`
	StartTime  time.Time `gorm:"type:timestamp;not null" json:"start_time"`
	EndTime    time.Time `gorm:"type:timestamp;not null" json:"end_time"`
	PeopleNum  uint8     `gorm:"not null" json:"people_num"`
	IsApproved bool      `gorm:"type:boolean;not null" json:"is_approved"`
}

func (r *CinemaReservation) TableName() string {
	return "cinema_reservations"
}

type GetReservationRequest struct {
	StartTimeFrom time.Time
	StartTimeTo   time.Time
	EndTimeTo     time.Time
	EndTimeFrom   time.Time
	PeopleNumFrom *uint8
	PeopleNumTo   *uint8
	UserID        uuid.UUID
	IsApproved    *bool
}
