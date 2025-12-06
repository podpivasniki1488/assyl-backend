package model

import (
	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/protopb"
)

type User struct {
	ID           uuid.UUID    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	FirstName    string       `gorm:"type:varchar;not null" json:"first_name"`
	LastName     string       `gorm:"type:varchar;not null" json:"last_name"`
	Username     string       `gorm:"type:varchar;not null;uniqueIndex:idx_user_username" json:"username"`
	UsernameType int          `gorm:"type:int;not null" default:"user" json:"username_type"`
	Password     string       `gorm:"type:varchar;not null" json:"password"`
	IsApproved   bool         `gorm:"type:boolean;not null" default:"false" json:"is_approved"`
	ApartmentID  uuid.UUID    `gorm:"type:uuid" json:"apartment_id"`
	RoleID       protopb.Role `gorm:"type:smallint;not null" json:"role_id"`
}
