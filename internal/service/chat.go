package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository"
	"github.com/podpivasniki1488/assyl-backend/pkg"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type chat struct {
	repo   *repository.Repository
	tracer trace.Tracer
	logger *slog.Logger
	hub    *pkg.Hub
}

func NewChat(repo *repository.Repository, logger *slog.Logger, hub *pkg.Hub) Chat {
	return &chat{
		repo:   repo,
		tracer: otel.Tracer("chatService"),
		logger: logger,
		hub:    hub,
	}
}

func (c *chat) SubscribeToUpdates(userID uuid.UUID) (chan pkg.HubMessage, func()) {
	_, span := c.tracer.Start(context.Background(), "chatService.SubscribeToUpdates")
	defer span.End()

	return c.hub.Subscribe(userID)
}

func (c *chat) CheckUserAllowance(ctx context.Context, userId, chatId uuid.UUID) (bool, error) {
	ctx, span := c.tracer.Start(ctx, "chatService.CheckUserAllowance")
	defer span.End()

	ok, err := c.repo.ChatRepo.CheckUserInChat(ctx, chatId, userId)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (c *chat) SendMessageToChat(ctx context.Context, message string, chatId, senderId uuid.UUID) error {
	ctx, span := c.tracer.Start(ctx, "chatService.SendMessageToChat")
	defer span.End()

	// check user allowance
	ok, err := c.repo.ChatRepo.CheckUserInChat(ctx, chatId, senderId)
	if err != nil {
		return err
	}

	if !ok {
		return model.ErrUserNotAllowed
	}

	// insert message to mongodb
	if _, err = c.repo.ChatRepo.InsertNewMessage(ctx, model.Message{
		ChatID:    chatId,
		SenderID:  senderId,
		Text:      message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		return err
	}

	// change LastMessagePreview, LastParticipantId, and LastMessageAt in goroutine
	go func() {
		if err = c.repo.ChatRepo.ChangeLastChatInfo(ctx, message, senderId, chatId); err != nil {
			c.logger.Error("could not update last chat info:", err)
		}
	}()

	go func() {
		if err = c.hub.Publish(senderId, pkg.HubMessage{
			ChatID: chatId,
			Text:   message,
		}); err != nil {
			c.logger.Error("could not publish message to hub:", err)
		}
	}()

	return nil
}

func (c *chat) StartNewChat(ctx context.Context, creatorID uuid.UUID, participantsIDs []uuid.UUID) (uuid.UUID, error) {
	ctx, span := c.tracer.Start(ctx, "chatService.StartChat")
	defer span.End()

	// TODO: maybe additionally check for chat existence?

	if len(participantsIDs) < 2 {
		return uuid.Nil, model.ErrSingleChatUser
	}

	seen := make(map[uuid.UUID]bool)

	for _, participant := range participantsIDs {
		if _, exists := seen[participant]; exists {
			return uuid.Nil, model.ErrCannotHaveDuplicates
		}

		seen[participant] = true
	}

	chatId, err := c.repo.ChatRepo.CreateChat(ctx, creatorID, participantsIDs)
	if err != nil {
		return uuid.Nil, err
	}

	return *chatId, nil
}
