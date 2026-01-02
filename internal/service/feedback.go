package service

import (
	"context"

	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository"
	"github.com/podpivasniki1488/assyl-backend/protopb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type feedback struct {
	repo   *repository.Repository
	tracer trace.Tracer
}

func NewFeedback(repo *repository.Repository) Feedback {
	return &feedback{repo, otel.Tracer("feedbackService")}
}

func (f *feedback) CreateFeedback(ctx context.Context, req model.Feedback) error {
	ctx, span := f.tracer.Start(ctx, "feedbackService.CreateFeedback")
	defer span.End()

	if err := f.repo.FeedbackRepo.CreateFeedback(ctx, &req); err != nil {
		return err
	}

	return nil
}

func (f *feedback) GetFeedbacks(ctx context.Context, req model.GetFeedbackRequest) ([]model.Feedback, error) {
	ctx, span := f.tracer.Start(ctx, "feedbackService.GetFeedbacks")
	defer span.End()

	resp, err := f.repo.FeedbackRepo.GetFeedbackByFilter(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (f *feedback) existType(enum protopb.FeedbackType) bool {
	_, ok := protopb.FeedbackType_value[enum.String()]

	return ok
}
