package slot

import (
	"context"
	"time"

	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type slotRepo struct {
	db *gorm.DB
}

func NewSlotRepo(db *gorm.DB) SlotRepo {
	return &slotRepo{db: db}
}

func (r *slotRepo) GetActiveTemplates(ctx context.Context) ([]model.SlotTemplate, error) {
	var t []model.SlotTemplate

	if err := r.db.WithContext(ctx).
		Where("is_active = true").
		Order("position asc").
		Find(&t).Error; err != nil {
		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	return t, nil
}

// EnsureDailySlots создает daily_slots на конкретную дату (если еще нет)
// tz — таймзона кинотеатра (например, Asia/Almaty)
func (r *slotRepo) EnsureDailySlots(ctx context.Context, date time.Time, tz *time.Location) error {
	templates, err := r.GetActiveTemplates(ctx)
	if err != nil {
		return err
	}

	// date (date-only): отрежем время
	y, m, d := date.In(tz).Date()
	dayStart := time.Date(y, m, d, 0, 0, 0, 0, tz)

	toTime := func(h, min int) time.Time {
		return time.Date(y, m, d, h, min, 0, 0, tz)
	}

	// ВАЖНО: start_time/end_time в templates храню как строки "10:00:00"
	// Чтобы быстро, просто сопоставим по position.
	// Если хочешь строго из БД time, можно парсить.
	makeRangeByPosition := func(pos int16) (time.Time, time.Time) {
		switch pos {
		case 1:
			return toTime(10, 0), toTime(12, 0)
		case 2:
			return toTime(14, 0), toTime(16, 0)
		case 3:
			return toTime(18, 0), toTime(20, 0)
		case 4:
			return toTime(22, 0), toTime(23, 30)
		default:
			return dayStart, dayStart
		}
	}

	slots := make([]model.DailySlot, 0, len(templates))
	for _, t := range templates {
		startAt, endAt := makeRangeByPosition(t.Position)
		slots = append(slots, model.DailySlot{
			SlotDate:   dayStart, // type: date в БД
			TemplateID: t.ID,
			Position:   t.Position,
			StartAt:    startAt.UTC(), // храним в UTC
			EndAt:      endAt.UTC(),
			IsEnabled:  true,
		})
	}

	if err = r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&slots).Error; err != nil {
		return model.ErrDBUnexpected.WithErr(err)
	}

	return nil
}

func (r *slotRepo) GetDailySlots(ctx context.Context, date time.Time, tz *time.Location) ([]model.DailySlot, error) {
	y, m, d := date.In(tz).Date()
	day := time.Date(y, m, d, 0, 0, 0, 0, tz)

	var ds []model.DailySlot
	if err := r.db.WithContext(ctx).
		Where("slot_date = ?", day).
		Order("position asc").
		Find(&ds).Error; err != nil {
		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	return ds, nil
}

// Free slots: те, где нет записи reservation_slots
func (r *slotRepo) GetFreeDailySlots(ctx context.Context, date time.Time, tz *time.Location) ([]model.DailySlot, error) {
	y, m, d := date.In(tz).Date()
	day := time.Date(y, m, d, 0, 0, 0, 0, tz)

	var ds []model.DailySlot
	if err := r.db.WithContext(ctx).
		Table("daily_slots ds").
		Select("ds.*").
		Joins("LEFT JOIN reservation_slots rs ON rs.daily_slot_id = ds.id").
		Where("ds.slot_date = ? AND ds.is_enabled = true AND rs.daily_slot_id IS NULL", day).
		Order("ds.position asc").
		Scan(&ds).Error; err != nil {
		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	return ds, nil
}

// Получить daily_slot_id по date + position
func (r *slotRepo) GetDailySlotIDsByPositions(ctx context.Context, date time.Time, positions []int16, tz *time.Location) (map[int16]uint64, error) {
	y, m, d := date.In(tz).Date()
	day := time.Date(y, m, d, 0, 0, 0, 0, tz)

	type row struct {
		Position int16
		ID       int64
	}

	var rows []row

	err := r.db.WithContext(ctx).
		Table("daily_slots").
		Select("position, id").
		Where("slot_date = ? AND position IN ?", day, positions).
		Scan(&rows).Error
	if err != nil {
		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	out := make(map[int16]uint64, len(rows))
	for _, r := range rows {
		out[r.Position] = uint64(r.ID)
	}
	return out, nil
}
