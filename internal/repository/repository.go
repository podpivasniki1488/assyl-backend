package repository

import (
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository/email"
	"github.com/podpivasniki1488/assyl-backend/internal/repository/user"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	UserRepo  user.UserRepo
	EmailRepo email.EmailRepo
}

func MustInitDb(dsn string) *gorm.DB {
	db, err := gorm.Open(
		postgres.Open(
			dsn,
		),
		&gorm.Config{},
	)
	if err != nil {
		panic(err)
	}

	if err = db.AutoMigrate(
		&model.User{},
		&model.Apartment{},
	); err != nil {
		panic(err)
	}

	return db
}

func NewRepository(db *gorm.DB, debug bool, trace trace.Tracer) *Repository {
	return &Repository{
		UserRepo:  user.NewUserRepository(db, debug, trace),
		EmailRepo: email.NewEmailRepo(trace),
	}
}
