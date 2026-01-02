package feedback

import (
	"context"

	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

type FeedbackRepo interface {
	CreateFeedback(ctx context.Context, feedback *model.Feedback) error
	GetFeedbackByFilter(ctx context.Context, req *model.GetFeedbackRequest) ([]model.Feedback, error)
}
