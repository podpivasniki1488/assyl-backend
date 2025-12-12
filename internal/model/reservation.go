package model

import (
	"time"

	"github.com/google/uuid"
)

type CinemaReservation struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;foreignKey:UserID;references:ID" json:"user_id"`
	From      time.Time `gorm:"type:timestamp;not null" json:"from"`
	To        time.Time `gorm:"type:timestamp;not null;check:from < to" json:"to"`
	PeopleNum uint8     `gorm:"not null" json:"people_num"`
}
