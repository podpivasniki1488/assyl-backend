package service

import (
	"context"
	"time"

	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type channelService struct {
	repo   *repository.Repository
	tracer trace.Tracer
}

func NewChannelService(repo *repository.Repository) Channel {
	return &channelService{repo, otel.Tracer("channelService")}
}

func (c *channelService) SendChannelMessage(ctx context.Context, msg model.ChannelMessage) error {
	ctx, span := c.tracer.Start(ctx, "channelService.SendChannelMessage")
	defer span.End()

	if err := c.repo.ChannelRepo.InsertNewMessage(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (c *channelService) GetByTimePeriod(ctx context.Context, from, to time.Time) ([]model.ChannelMessage, error) {
	ctx, span := c.tracer.Start(ctx, "channelService.GetByTimePeriod")
	defer span.End()

	if from.After(to) {
		return nil, nil
	}

	res, err := c.repo.ChannelRepo.GetMessageByTime(ctx, from, to)
	if err != nil {
		return nil, err
	}

	return res, nil
}
