package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository"
	"github.com/podpivasniki1488/assyl-backend/protopb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type userManagement struct {
	repo   *repository.Repository
	tracer trace.Tracer
}

func NewUserManagement(repo *repository.Repository) UserManagement {
	return &userManagement{repo, otel.Tracer("userManagement")}
}

func (u *userManagement) DeleteUserByUsername(ctx context.Context, username string) error {
	userToDelete, err := u.repo.UserRepo.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	if userToDelete.RoleID == protopb.Role_GOD || userToDelete.RoleID == protopb.Role_ADMIN {
		return model.ErrAdminsCannotBeDeleted
	}

	if err = u.repo.UserRepo.DeleteByUsername(ctx, username); err != nil {
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

	if err = u.checkUserAccess(ctx, *currUser, *ap); err != nil {
		return err
	}

	ap.OwnerId = &currUser.ID

	if err = u.repo.ApartmentRepo.UpdateApartment(ctx, ap); err != nil {
		return err
	}

	return nil
}

func (u *userManagement) checkUserAccess(ctx context.Context, currUser model.User, apartment model.Apartment) error {
	ctx, span := u.tracer.Start(ctx, "userManagement_checkUserAccess")
	defer span.End()

	if currUser.RoleID == protopb.Role_ADMIN || currUser.RoleID == protopb.Role_GOD {
		return nil
	}

	if apartment.OwnerId != nil && currUser.ID != *apartment.OwnerId {
		return model.ErrApartmentAlreadyBound
	}

	return nil
}
