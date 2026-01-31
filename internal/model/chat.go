package model

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Chat struct {
	ID                 uuid.UUID  `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	IsPrivate          bool       `gorm:"type:boolean;not null" json:"is_private"`
	LastMessageAt      *time.Time `gorm:"type:timestamp" json:"last_message_at"`
	LastMessagePreview string     `gorm:"type:varchar" json:"last_message_preview"`
	LastParticipantId  uuid.UUID  `gorm:"type:uuid" json:"last_participant_id"`
}

func (Chat) TableName() string {
	return "chats"
}

type ChatParticipant struct {
	ID       uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ChatID   uuid.UUID `gorm:"type:uuid;not null" json:"chat_id"`
	JoinedAt time.Time `gorm:"type:timestamp;not null" json:"joined_at"`
}

func (ChatParticipant) TableName() string {
	return "chat_participants"
}

type Message struct {
	ID               bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ChatID           uuid.UUID     `bson:"chat_id" json:"chat_id"`
	Text             string        `bson:"text" json:"text"`
	SenderID         uuid.UUID     `bson:"sender_id" json:"sender_id"`
	AttachmentLinks  []string      `bson:"attachment_links" json:"attachment_links"`
	CreatedAt        time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time     `bson:"updated_at" json:"updated_at"`
	DeletedAt        *time.Time    `bson:"deleted_at,omitempty" json:"deleted_at"`
	ReplyToMessageID *uuid.UUID    `bson:"reply_to_message_id,omitempty" json:"reply_to_message_id"`
}
