package repository

import (
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository/apartment"
	"github.com/podpivasniki1488/assyl-backend/internal/repository/channel"
	"github.com/podpivasniki1488/assyl-backend/internal/repository/email"
	"github.com/podpivasniki1488/assyl-backend/internal/repository/reservation"
	"github.com/podpivasniki1488/assyl-backend/internal/repository/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	UserRepo        user.UserRepo
	EmailRepo       email.EmailRepo
	ApartmentRepo   apartment.ApartmentRepo
	ReservationRepo reservation.ReservationRepo
	ChannelRepo     channel.ChanRepo
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
		&model.CinemaReservation{},
		&model.ChannelMessage{},
	); err != nil {
		panic(err)
	}

	return db
}

func NewRepository(db *gorm.DB, debug bool, gmailUsername, gmailPsw string) *Repository {
	return &Repository{
		UserRepo:        user.NewUserRepository(db, debug),
		EmailRepo:       email.NewEmailRepo(gmailUsername, gmailPsw),
		ApartmentRepo:   apartment.NewApartmentRepo(db),
		ReservationRepo: reservation.NewReservationRepository(db),
		ChannelRepo:     channel.NewChanRepository(db),
	}
}
