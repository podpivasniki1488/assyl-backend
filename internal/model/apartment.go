package model

import "github.com/google/uuid"

type Apartment struct {
	Id         uuid.UUID `gorm:"type:uuid;not null;default:gen_random_uuid()" json:"id"`
	Floor      uint8     `gorm:"not null;uniqueIndex:idx_floor_door" json:"floor"`
	DoorNumber uint16    `gorm:"not null;uniqueIndex:idx_floor_door" json:"door_number"`
}
