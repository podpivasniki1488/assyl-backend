package service

import (
	"context"
	"regexp"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository"
	"github.com/podpivasniki1488/assyl-backend/protopb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type reservation struct {
	tracer trace.Tracer
	repo   *repository.Repository
}

var phoneNumRegex = regexp.MustCompile(`^\+\d{11}$`)

func NewReservation(repo *repository.Repository) Reservation {
	return &reservation{
		tracer: otel.Tracer("reservationService"),
		repo:   repo,
	}
}

func (r *reservation) GetUserReservations(ctx context.Context, req model.CinemaReservation) ([]model.CinemaReservation, error) {
	ctx, span := r.tracer.Start(ctx, "reservation.GetReservation")
	defer span.End()

	user, err := r.repo.UserRepo.FindById(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	isApproved := true

	reservations, err := r.repo.ReservationRepo.GetByFilters(ctx, &model.GetReservationRequest{
		StartTimeFrom: req.StartTime,
		EndTimeTo:     req.EndTime,
		IsApproved:    &isApproved,
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

func (r *reservation) MakeReservation(ctx context.Context, req *model.CinemaReservation, role, username string) error {
	ctx, span := r.tracer.Start(ctx, "reservation.MakeReservation")
	defer span.End()

	if req.PeopleNum > 12 {
		return model.ErrTooManyPeople
	}

	if req.StartTime.After(req.EndTime) {
		return model.ErrInvalidInput
	}

	if req.EndTime.Sub(req.StartTime) > 12*time.Hour {
		return model.ErrReservationImpossible
	}

	reservations, err := r.repo.ReservationRepo.GetByFilters(ctx, &model.GetReservationRequest{
		StartTimeFrom: req.StartTime,
		EndTimeTo:     req.EndTime,
	})
	if err != nil {
		return err
	}

	if len(reservations) != 0 {
		return model.ErrCinemaBusy
	}

	req.IsApproved = false
	if role == protopb.Role_ADMIN.String() || role == protopb.Role_GOD.String() {
		req.IsApproved = true
	}

	if phoneNumRegex.MatchString(username) {
		req.PhoneNum = username
	}

	if err = r.repo.ReservationRepo.CreateReservation(ctx, req); err != nil {
		return err
	}

	return nil
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

func (r *reservation) GetFreeSlots(ctx context.Context, from, to time.Time) ([]model.TimeRange, error) {
	ctx, span := r.tracer.Start(ctx, "reservation.GetFreeSlots")
	defer span.End()

	if !from.Before(to) {
		return []model.TimeRange{}, nil
	}

	// TODO: нужно ли нам сортировать по approved?
	reservations, err := r.repo.ReservationRepo.GetByFilters(ctx, &model.GetReservationRequest{
		StartTimeFrom: from,
		EndTimeTo:     to,
	})
	if err != nil {
		return nil, err
	}

	if len(reservations) == 0 {
		return []model.TimeRange{{Start: from, End: to}}, nil
	}

	busy := make([]model.TimeRange, 0, len(reservations))
	for _, res := range reservations {
		s := res.StartTime
		e := res.EndTime

		// на всякий случай защита от мусора
		if !s.Before(e) {
			continue
		}

		// clamp to [from,to]
		if s.Before(from) {
			s = from
		}
		if e.After(to) {
			e = to
		}

		// после clamp может стать пустым
		if !s.Before(e) {
			continue
		}

		busy = append(busy, model.TimeRange{Start: s, End: e})
	}

	if len(busy) == 0 {
		return []model.TimeRange{{Start: from, End: to}}, nil
	}

	sort.Slice(busy, func(i, j int) bool {
		if busy[i].Start.Equal(busy[j].Start) {
			return busy[i].End.Before(busy[j].End)
		}
		return busy[i].Start.Before(busy[j].Start)
	})

	merged := make([]model.TimeRange, 0, len(busy))
	cur := busy[0]
	for i := 1; i < len(busy); i++ {
		n := busy[i]

		// если следующий начинается ДО или РОВНО в конец текущего — сливаем
		// (ровно = "соприкосновение" без свободной щели)
		if !n.Start.After(cur.End) {
			if n.End.After(cur.End) {
				cur.End = n.End
			}
			continue
		}

		merged = append(merged, cur)
		cur = n
	}
	merged = append(merged, cur)

	// 4) Вычисляем комплемент: free = [from,to] \ merged
	free := make([]model.TimeRange, 0, len(merged)+1)

	// до первой busy
	if from.Before(merged[0].Start) {
		free = append(free, model.TimeRange{Start: from, End: merged[0].Start})
	}

	// между busy
	for i := 0; i < len(merged)-1; i++ {
		a := merged[i]
		b := merged[i+1]
		if a.End.Before(b.Start) {
			free = append(free, model.TimeRange{Start: a.End, End: b.Start})
		}
	}

	// после последней busy
	last := merged[len(merged)-1]
	if last.End.Before(to) {
		free = append(free, model.TimeRange{Start: last.End, End: to})
	}

	return free, nil
}
