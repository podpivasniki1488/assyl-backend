package model

import "github.com/google/uuid"

type Apartment struct {
	Id         uuid.UUID `gorm:"type:uuid;not null;default:gen_random_uuid()"`
	Floor      uint8     `gorm:"not null"`
	DoorNumber uint16    `gorm:"not null"`
}
