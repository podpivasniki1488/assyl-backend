package feedback

import (
	"context"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/protopb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type feedbackRepository struct {
	db     *gorm.DB
	tracer trace.Tracer
	debug  bool
}

func NewFeedbackRepository(db *gorm.DB) FeedbackRepo {
	return &feedbackRepository{db, otel.Tracer("feedbackRepo"), true}
}

func (f *feedbackRepository) CreateFeedback(ctx context.Context, feedback *model.Feedback) error {
	ctx, span := f.tracer.Start(ctx, "feedbackRepo.CreateFeedback")
	defer span.End()

	query := f.db.WithContext(ctx).Model(&feedback)

	if f.debug {
		query = query.Debug()
	}

	if err := query.Create(feedback).Error; err != nil {
		return model.ErrDBUnexpected.WithErr(err)
	}

	return nil
}

func (f *feedbackRepository) GetFeedbackByFilter(ctx context.Context, req *model.GetFeedbackRequest) ([]model.Feedback, error) {
	ctx, span := f.tracer.Start(ctx, "feedbackRepo.GetFeedbackByFilter")
	defer span.End()

	query := f.db.WithContext(ctx).Model(&model.Feedback{})

	if !req.CreatedAtFrom.IsZero() {
		query = query.Where("created_at >= ?", req.CreatedAtFrom)
	}

	if !req.CreatedAtTo.IsZero() {
		query = query.Where("created_at <= ?", req.CreatedAtTo)
	}

	if req.UserID != uuid.Nil {
		query = query.Where("user_id = ?", req.UserID)
	}

	if req.ID != uuid.Nil {
		query = query.Where("id = ?", req.ID)
	}

	if val, ok := protopb.FeedbackType_value[req.FeedbackType]; ok {
		query = query.Where("feedback_type = ?", val)
	}

	if f.debug {
		query = query.Debug()
	}

	var resp []model.Feedback
	if err := query.Find(&resp).Error; err != nil {
		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	return resp, nil
}
