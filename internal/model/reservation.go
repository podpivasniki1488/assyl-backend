package model

import (
	"time"

	"github.com/google/uuid"
)

type CinemaReservation struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	From      time.Time `gorm:"type:timestamp;not null"`
	To        time.Time `gorm:"type:timestamp;not null"`
	PeopleNum uint8     `gorm:"not null"`
}
