package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type userRepository struct {
	db     *gorm.DB
	logger slog.Logger
	tracer trace.Tracer
	debug  bool
}

func NewUserRepository(db *gorm.DB, debug bool) UserRepo {
	return &userRepository{
		db:     db,
		debug:  debug,
		tracer: otel.Tracer("userRepository"),
	}
}

func (u *userRepository) FindById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	ctx, span := u.tracer.Start(ctx, "userRepository.FindById")
	defer span.End()

	var user model.User
	query := u.db.
		WithContext(ctx).
		Where("id = ?", id.String())

	if u.debug {
		query = query.Debug()
	}

	if err := query.Find(&user).Error; err != nil {
		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	return &user, nil
}

func (u *userRepository) FindByApartmentId(ctx context.Context, apartmentId uuid.UUID) ([]model.User, error) {
	ctx, span := u.tracer.Start(ctx, "userRepository.FindByApartmentId")
	defer span.End()

	var users []model.User
	query := u.db.
		WithContext(ctx).
		Where("apartment_id = ?", apartmentId.String()).
		Find(&users)

	if u.debug {
		query = query.Debug()
	}

	if err := query.Error; err != nil {
		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	return users, nil
}

func (u *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	ctx, span := u.tracer.Start(ctx, "userRepository.FindByusername")
	defer span.End()

	var user model.User
	query := u.db.
		Model(&user).
		WithContext(ctx).
		Where("username = ?", username)

	if u.debug {
		query = query.Debug()
	}

	if err := query.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}

		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	return &user, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	ctx, span := u.tracer.Start(ctx, "userRepository.UpdateUser")
	defer span.End()

	var findResp model.User
	if err := u.db.WithContext(ctx).
		Where("username = ?", user.Username).
		First(&findResp).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}

		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	user.ID = findResp.ID

	if err := u.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	var updateResp model.User
	if err := u.db.
		WithContext(ctx).
		First(&updateResp, "id = ?", user.ID).
		Error; err != nil {
		return nil, model.ErrDBUnexpected.WithErr(err)
	}

	return &updateResp, nil
}

func (u *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	ctx, span := u.tracer.Start(ctx, "userRepository.CreateUser")
	defer span.End()

	var predictedUser []model.User
	findQuery := u.db.WithContext(ctx).
		Where("username = ?", user.Username).
		Find(&predictedUser)

	if findQuery.Error != nil {
		return model.ErrDBUnexpected.WithErr(findQuery.Error)
	}

	if len(predictedUser) > 0 {
		return model.ErrUserAlreadyExists
	}

	insertQuery := u.db.WithContext(ctx).Create(user)
	if err := insertQuery.Error; err != nil {
		return model.ErrDBUnexpected.WithErr(err)
	}

	return nil
}

func (u *userRepository) DeleteByUsername(ctx context.Context, username string) error {
	ctx, span := u.tracer.Start(ctx, "userRepository.DeleteByUsername")
	defer span.End()

	var user model.User
	query := u.db.
		WithContext(ctx).
		Where("username = ?", username)

	if u.debug {
		query = query.Debug()
	}

	if err := query.Delete(&user).Error; err != nil {
		return model.ErrDBUnexpected.WithErr(err)
	}

	return nil
}
