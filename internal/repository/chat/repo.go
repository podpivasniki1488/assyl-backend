package chat

import (
	"context"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

type ChatRepo interface {
	InsertNewMessage(ctx context.Context, msg model.Message) (string, error)
	CheckUserInChat(ctx context.Context, userID, chatID uuid.UUID) (bool, error)
	ChangeLastChatInfo(ctx context.Context, msg string, userID, chatID uuid.UUID) error
	CreateChat(ctx context.Context, creatorID uuid.UUID, participantsIDs []uuid.UUID) (*uuid.UUID, error)
}
