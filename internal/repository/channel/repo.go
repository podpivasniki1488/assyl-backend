package channel

import (
	"context"
	"time"

	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

type ChanRepo interface {
	InsertNewMessage(ctx context.Context, msg model.ChannelMessage) error
	GetMessageByTime(ctx context.Context, from, to time.Time) ([]model.ChannelMessage, error)
}
