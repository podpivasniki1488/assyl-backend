package channel

import (
	"context"
	"time"

	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type chanRepo struct {
	db     *gorm.DB
	debug  bool
	tracer trace.Tracer
}

func NewChanRepository(db *gorm.DB) ChanRepo {
	return &chanRepo{db, true, otel.Tracer("channelRepo")}
}

func (c *chanRepo) InsertNewMessage(ctx context.Context, msg model.ChannelMessage) error {
	ctx, span := c.tracer.Start(ctx, "channelRepo.InsertNewMessage")
	defer span.End()

	query := c.db.WithContext(ctx)

	if c.debug {
		query = query.Debug()
	}

	if err := query.Create(&msg).Error; err != nil {
		return model.ErrDBUnexpected.WithErr(err)
	}

	return nil
}

func (c *chanRepo) GetMessageByTime(ctx context.Context, from, to time.Time) ([]model.ChannelMessage, error) {
	ctx, span := c.tracer.Start(ctx, "channelRepo.GetMessageByTime")
	defer span.End()

	query := c.db.WithContext(ctx).Where("created_at BETWEEN ? AND ?", from, to)

	if c.debug {
		query = query.Debug()
	}

	var res []model.ChannelMessage
	if err := query.Find(&res).Error; err != nil {
		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	return res, nil
}
