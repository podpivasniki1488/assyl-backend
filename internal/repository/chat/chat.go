package chat

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type chat struct {
	pgDB        *gorm.DB
	mongoClient *mongo.Client
	tracer      trace.Tracer
	debug       bool
}

func NewChatRepo(pgDB *gorm.DB, mongoClient *mongo.Client, debug bool) ChatRepo {
	return &chat{pgDB, mongoClient, otel.Tracer("chatRepo"), debug}
}

func (c *chat) InsertNewMessage(ctx context.Context, msg model.Message) (string, error) {
	ctx, span := c.tracer.Start(ctx, "chatRepo.InsertNewMessage")
	defer span.End()

	coll := c.mongoClient.
		Database("assyl").
		Collection("messages")

	res, err := coll.InsertOne(ctx, msg)
	if err != nil {
		return "", err
	}

	id, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return "", model.ErrDBUnexpected.WithErr(errors.New("could not cast id to string"))
	}

	return id.Hex(), nil
}

func (c *chat) CheckUserInChat(ctx context.Context, userID, chatID uuid.UUID) (bool, error) {
	ctx, span := c.tracer.Start(ctx, "chatRepo.CheckUserInChat")
	defer span.End()

	var res model.ChatParticipant
	query := c.pgDB.Model(&model.ChatParticipant{}).
		Where("user_id = ? AND chat_id = ?", userID, chatID).
		First(&res)

	if err := query.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, model.ErrDBUnexpected.WithErr(err)
	}

	return true, nil
}

func (c *chat) ChangeLastChatInfo(ctx context.Context, msg string, userID, chatID uuid.UUID) error {
	ctx, span := c.tracer.Start(ctx, "chatRepo.ChangeLastChatInfo")
	defer span.End()

	getChatQuery := c.pgDB.Model(&model.Chat{}).
		First(&model.Chat{}, "id = ?", chatID)

	if err := getChatQuery.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrChatNotFound
		}

		return model.ErrDBUnexpected.WithErr(err)
	}

	updateChatQuery := c.pgDB.Model(&model.Chat{}).
		Where("id = ?", chatID).Updates(map[string]interface{}{
		"last_message_at":      time.Now(),
		"last_message_preview": msg,
		"last_participant_id":  userID,
	})

	if err := updateChatQuery.Error; err != nil {
		return model.ErrDBUnexpected.WithErr(err)
	}

	return nil
}

func (c *chat) CreateChat(ctx context.Context, creatorID uuid.UUID, participantsIDs []uuid.UUID) (*uuid.UUID, error) {
	ctx, span := c.tracer.Start(ctx, "chatRepo.CreateChat")
	defer span.End()

	tx := c.pgDB.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return nil, model.ErrDBUnexpected.WithErr(err)
	}
	defer tx.Rollback()

	isPrivate := false
	if len(participantsIDs) == 2 {
		isPrivate = true
	}

	newChat := model.Chat{
		IsPrivate:          isPrivate,
		LastMessageAt:      nil,
		LastMessagePreview: "",
		LastParticipantId:  creatorID,
	}

	createChatQuery := tx.Model(&model.Chat{}).Create(&newChat)

	if err := createChatQuery.Error; err != nil {
		tx.Rollback()
		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	for _, participant := range participantsIDs {
		createParticipant := model.ChatParticipant{
			UserID:   participant,
			ChatID:   newChat.ID,
			JoinedAt: time.Now(),
		}

		if err := tx.Model(&model.ChatParticipant{}).Create(&createParticipant).Error; err != nil {
			tx.Rollback()
			return nil, model.ErrDBUnexpected.WithErr(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	return &newChat.ID, nil
}
