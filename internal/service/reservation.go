package service

import (
	"context"
	"errors"
	"regexp"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository"
	"github.com/podpivasniki1488/assyl-backend/protopb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type reservation struct {
	tracer trace.Tracer
	repo   *repository.Repository
}

var phoneNumRegex = regexp.MustCompile(`^\+\d{11}$`)

const totalFreeReservations = 5

func NewReservation(repo *repository.Repository) Reservation {
	return &reservation{
		tracer: otel.Tracer("reservationService"),
		repo:   repo,
	}
}

func (r *reservation) GetUserReservations(ctx context.Context, userId uuid.UUID, start, end time.Time) ([]model.CinemaReservation, error) {
	ctx, span := r.tracer.Start(ctx, "reservation.GetReservation")
	defer span.End()

	startTime := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	endTime := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)

	user, err := r.repo.UserRepo.FindById(ctx, userId)
	if err != nil {
		return nil, err
	}

	reservations, err := r.repo.ReservationRepo.GetByFilters(ctx, &model.GetReservationRequest{
		StartTimeFrom: startTime,
		EndTimeTo:     endTime,
	})
	if err != nil {
		return nil, err
	}

	return r.filterReservation(ctx, reservations, *user)
}

func (r *reservation) GetUnfilteredReservations(ctx context.Context, req model.CinemaReservation) ([]model.CinemaReservation, error) {
	ctx, span := r.tracer.Start(ctx, "reservation.GetReservation")
	defer span.End()

	reservations, err := r.repo.ReservationRepo.GetByFilters(ctx, &model.GetReservationRequest{
		StartTimeFrom: req.StartTime,
		EndTimeTo:     req.EndTime,
	})
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *reservation) filterReservation(ctx context.Context, reservations []model.CinemaReservation, user model.User) ([]model.CinemaReservation, error) {
	ctx, span := r.tracer.Start(ctx, "reservation.filterReservation")
	defer span.End()

	if user.RoleID == protopb.Role_ADMIN || user.RoleID == protopb.Role_GOD {
		return reservations, nil
	}

	res := make([]model.CinemaReservation, 0, len(reservations))
	for _, resv := range reservations {
		if resv.UserID == user.ID {
			res = append(res, resv)
		}
	}

	return res, nil
}

func (r *reservation) MakeReservation(
	ctx context.Context,
	userID uuid.UUID,
	date time.Time,
	positions []int16,
	peopleNum uint8,
	role, username string,
) (reservationLeft int, err error) {
	ctx, span := r.tracer.Start(ctx, "reservation.MakeReservation")
	defer span.End()

	if peopleNum > 12 {
		return 0, model.ErrTooManyPeople
	}

	if len(positions) != 1 && len(positions) != 2 {
		return 0, model.ErrInvalidInput
	}

	sort.Slice(positions, func(i, j int) bool {
		return positions[i] < positions[j]
	})

	if len(positions) == 2 && positions[1] != positions[0]+1 {
		return 0, model.ErrInvalidInput
	}

	if err = r.repo.SlotRepo.EnsureDailySlots(ctx, date, time.UTC); err != nil {
		return 0, err
	}

	idMap, err := r.repo.SlotRepo.GetDailySlotIDsByPositions(ctx, date, positions, time.UTC)
	if err != nil {
		return 0, err
	}

	dailyIDs := make([]uint64, 0, len(positions))
	for _, p := range positions {
		id, ok := idMap[p]
		if !ok {
			return 0, model.ErrInvalidInput
		}
		dailyIDs = append(dailyIDs, id)
	}

	dSlots, err := r.repo.SlotRepo.GetDailySlots(ctx, date, time.UTC)
	if err != nil {
		return 0, err
	}

	var start, end time.Time
	posToSlot := map[int16]model.DailySlot{}
	for _, s := range dSlots {
		posToSlot[s.Position] = s
	}

	start = posToSlot[positions[0]].StartAt
	end = posToSlot[positions[len(positions)-1]].EndAt

	res := &model.CinemaReservation{
		ID:         uuid.New(),
		UserID:     userID,
		StartTime:  start,
		EndTime:    end,
		PeopleNum:  peopleNum,
		IsApproved: false,
		PhoneNum:   "",
	}

	if role == protopb.Role_ADMIN.String() || role == protopb.Role_GOD.String() {
		res.IsApproved = true
	}

	if phoneNumRegex.MatchString(username) {
		res.PhoneNum = username
	}

	user, err := r.repo.UserRepo.FindById(ctx, userID)
	if err != nil {
		return 0, err
	}

	userFromSameApart, err := r.repo.UserRepo.FindByApartmentId(ctx, user.ApartmentID)
	if err != nil {
		return 0, err
	}

	reservationLeft = totalFreeReservations
	for _, u := range userFromSameApart {
		reservations, err := r.repo.ReservationRepo.GetByFilters(ctx, &model.GetReservationRequest{
			StartTimeFrom: start,
			EndTimeTo:     end,
			UserID:        u.ID,
		})
		if err != nil {
			return 0, err
		}

		reservationLeft -= len(reservations)
	}

	if err = r.repo.ReservationRepo.CreateReservationsWithSlots(ctx, res, dailyIDs); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, model.ErrCinemaBusy
		}

		return 0, err
	}

	return reservationLeft, nil
}

func (r *reservation) ApproveReservation(ctx context.Context, id uuid.UUID) error {
	ctx, span := r.tracer.Start(ctx, "reservation.ApproveReservation")
	defer span.End()

	res, err := r.repo.ReservationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if res.IsApproved {
		return nil
	}

	res.IsApproved = true

	if err = r.repo.ReservationRepo.ApproveReservation(ctx, id); err != nil {
		return err
	}

	return nil
}

func (r *reservation) GetFreeSlots(ctx context.Context, date time.Time) ([]model.DailySlot, [][2]int16, error) {
	ctx, span := r.tracer.Start(ctx, "reservation.GetFreeSlots")
	defer span.End()

	if err := r.repo.SlotRepo.EnsureDailySlots(ctx, date, time.UTC); err != nil {
		return nil, nil, err
	}

	free, err := r.repo.SlotRepo.GetFreeDailySlots(ctx, date, time.UTC)
	if err != nil {
		return nil, nil, err
	}

	freePos := make(map[int16]bool, len(free))
	for _, s := range free {
		freePos[s.Position] = true
	}

	pairs := make([][2]int16, 0, 4)
	for _, s := range free {
		if freePos[s.Position+1] {
			pairs = append(pairs, [2]int16{s.Position, s.Position + 1})
		}
	}

	return free, pairs, nil
}
