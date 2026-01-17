package slot

import (
	"context"
	"time"

	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

type SlotRepo interface {
	GetActiveTemplates(ctx context.Context) ([]model.SlotTemplate, error)
	EnsureDailySlots(ctx context.Context, date time.Time, tz *time.Location) error
	GetDailySlots(ctx context.Context, date time.Time, tz *time.Location) ([]model.DailySlot, error)
	GetFreeDailySlots(ctx context.Context, date time.Time, tz *time.Location) ([]model.DailySlot, error)
	GetDailySlotIDsByPositions(ctx context.Context, date time.Time, positions []int16, tz *time.Location) (map[int16]uint64, error)
}
