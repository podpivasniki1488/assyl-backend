package model

import (
	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/protopb"
)

type User struct {
	ID           uuid.UUID    `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	FirstName    string       `gorm:"type:varchar;not null"`
	LastName     string       `gorm:"type:varchar;not null"`
	Username     string       `gorm:"type:varchar;not null;uniqueIndex:idx_user_username"`
	UsernameType int          `gorm:"type:int;not null" default:"user"`
	Password     string       `gorm:"type:varchar;not null"`
	IsApproved   bool         `gorm:"type:boolean;not null" default:"false"`
	ApartmentID  uuid.UUID    `gorm:"type:uuid"` //TODO: create apartment table
	RoleID       protopb.Role `gorm:"type:smallint;not null"`
}
