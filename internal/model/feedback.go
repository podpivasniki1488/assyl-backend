package model

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	Id        uuid.UUID `gorm:"type:uuid;not null;default:gen_random_uuid()" json:"id"`
	UserId    uuid.UUID `gorm:"type:uuid;not null" json:"author_id"`
	Text      string    `gorm:"type:varchar;not null" json:"text"`
	CreatedAt time.Time `gorm:"type:timestamp;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null" json:"updated_at"`
}
