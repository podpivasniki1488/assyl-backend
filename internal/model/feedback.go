package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/protopb"
)

type Feedback struct {
	CreatedAt    time.Time            `gorm:"type:timestamp;not null" json:"created_at"`
	UpdatedAt    time.Time            `gorm:"type:timestamp;not null" json:"updated_at"`
	Id           uuid.UUID            `gorm:"type:uuid;not null;default:gen_random_uuid()" json:"id"`
	UserId       uuid.UUID            `gorm:"type:uuid;not null" json:"author_id"`
	Text         string               `gorm:"type:varchar;not null" json:"text"`
	FeedbackType protopb.FeedbackType `gorm:"type:smallint;not null" json:"feedback_type"`
}

func (Feedback) TableName() string {
	return "feedbacks"
}

type GetFeedbackRequest struct {
	CreatedAtFrom time.Time
	CreatedAtTo   time.Time
	ID            uuid.UUID
	UserID        uuid.UUID
	FeedbackType  string
	Text          string
}
