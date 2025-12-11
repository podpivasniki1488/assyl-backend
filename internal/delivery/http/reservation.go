package http

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

func (h *httpDelivery) registerReservationHandlers(v1 *echo.Group) {
	reservation := v1.Group("/reservation")

	reservation.Use(h.registerJWTMiddleware())
	reservation.POST("", h.createReservation, h.getJWTData())
	reservation.GET("", h.getReservation)
}

func (h *httpDelivery) getReservation(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.getReservation")
	defer span.End()

	var req getReservationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := h.service.Reservation.GetUnfilteredReservations(ctx, model.CinemaReservation{
		From:      *req.From,
		To:        *req.To,
		PeopleNum: *req.PeopleNum,
	})
	if err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusOK, DefaultResponse[[]model.CinemaReservation]{
		Status: "success",
		Data:   resp,
	})
}

func (h *httpDelivery) createReservation(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.createReservation")
	defer span.End()

	var req createReservationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userId, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "failed to get user_id from context")
	}

	if err := h.service.Reservation.MakeReservation(ctx, &model.CinemaReservation{
		From:      req.From,
		To:        req.To,
		PeopleNum: req.PeopleNum,
		UserID:    userId,
	}); err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusCreated, DefaultResponse[string]{
		Status: "ok",
		Data:   "",
	})
}

type getReservationRequest struct {
	From      *time.Time `query:"from"`
	To        *time.Time `query:"to"`
	PeopleNum *uint8     `query:"people_num"`
}

type createReservationRequest struct {
	From      time.Time `json:"from" validate:"required"`
	To        time.Time `json:"to" validate:"required,gtfield=From"`
	PeopleNum uint8     `json:"people_num" validate:"required,gt=1,lt=12"`
}
