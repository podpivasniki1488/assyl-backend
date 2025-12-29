package http

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

// TODO: here we need to add user handlers and first handler must be giving users cinema reservations

func (h *httpDelivery) registerUserHandlers(v1 *echo.Group) {
	user := v1.Group("/user")

	reservation := user.Group("/reservation")
	reservation.POST("/cinema", h.reserveCinema, h.getJWTData())

}

func (h *httpDelivery) reserveCinema(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.reserveCinema")
	defer span.End()

	var req reserveCinemaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	userId, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("user_id not found in context"))
	}

	id, err := uuid.Parse(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
	}

	if err = h.service.Reservation.MakeReservation(ctx, &model.CinemaReservation{
		UserID:    id,
		StartTime: req.From,
		EndTime:   req.To,
		PeopleNum: req.PeopleNum,
	}); err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusCreated, DefaultResponse[string]{
		Status: "success",
		Data:   "",
	})
}

type reserveCinemaRequest struct {
	From      time.Time `json:"from" validate:"required"`
	To        time.Time `json:"to" validate:"required"`
	PeopleNum uint8     `json:"people_num" validate:"required"`
}
