package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/repository"
	"go.opentelemetry.io/otel/trace"
)

type userManagement struct {
	repo   *repository.Repository
	tracer trace.Tracer
}

func NewUserManagement(repo *repository.Repository, tracer trace.Tracer) UserManagement {
	return &userManagement{repo, tracer}
}

func (u *userManagement) DeleteUserByEmail(ctx context.Context, email string) error {
	if err := u.repo.UserRepo.DeleteByUsername(ctx, email); err != nil {
		return err
	}

	return nil
}

func (u *userManagement) BindApartmentToUser(ctx context.Context, username string, apartmentId uuid.UUID) error {
	ctx, span := u.tracer.Start(ctx, "userManagement_BindApartmentToUser")
	defer span.End()

	ap, err := u.repo.ApartmentRepo.GetApartmentByID(ctx, apartmentId)
	if err != nil {
		return err
	}

	currUser, err := u.repo.UserRepo.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	ap.OwnerId = currUser.ID

	if err = u.repo.ApartmentRepo.UpdateApartment(ctx, ap); err != nil {
		return err
	}

	return nil
}
