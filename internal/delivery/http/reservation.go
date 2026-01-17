package http

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/protopb"
)

func (h *httpDelivery) registerReservationHandlers(v1 *echo.Group) {
	reservation := v1.Group("/reservation")

	reservation.Use(h.registerJWTMiddleware())
	reservation.POST("", h.createReservation, h.getJWTData())
	reservation.GET("", h.getReservation, h.getJWTData())
	reservation.PATCH("/approve", h.approveReservation, h.getJWTData())
	reservation.GET("/free-slots", h.getFreeSlots)
}

// getReservation godoc
//
//	@Summary		Get user reservations
//	@Description	Возвращает список бронирований текущего пользователя за период (from-to).
//	@Tags			reservation
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			from	query		string										true	"Start datetime (RFC3339)"	example(2025-12-29T10:00:00Z)
//	@Param			to		query		string										true	"End datetime (RFC3339)"	example(2025-12-29T12:00:00Z)
//	@Success		200		{object}	DefaultResponse[[]model.CinemaReservation]	"Успех"
//	@Failure		400		{object}	DefaultResponse[error]						"Невалидный запрос"
//	@Failure		401		{object}	DefaultResponse[error]						"Неавторизован"
//	@Failure		403		{object}	DefaultResponse[error]						"Нет доступа"
//	@Failure		500		{object}	DefaultResponse[error]						"Внутренняя ошибка сервера"
//	@Router			/reservation [get]
func (h *httpDelivery) getReservation(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.getReservation")
	defer span.End()

	userId, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("failed to get user_id from context"))
	}

	parsed, err := uuid.Parse(userId)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("invalid user id"))
	}

	var req getReservationRequest
	if err = c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid date format (YYYY-MM-DD)"))
	}

	resp, err := h.service.Reservation.GetUserReservations(ctx, model.CinemaReservation{
		UserID: parsed,
	}, parsedDate)
	if err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusOK, DefaultResponse[[]model.CinemaReservation]{
		Status: "success",
		Data:   resp,
	})
}

// approveReservation godoc
//
//	@Summary		Approve reservation
//	@Description	Одобряет бронирование по reservation_id. Доступно только администратору.
//	@Tags			reservation
//	@Security		BearerAuth
//	@Produce		json
//	@Param			reservation_id	query	string	true	"reservation_id"
//	@Success		204				"Одобрено"
//	@Failure		400				{object}	DefaultResponse[error]	"Невалидный запрос"
//	@Failure		401				{object}	DefaultResponse[error]	"Неавторизован"
//	@Failure		403				{object}	DefaultResponse[error]	"Только админ может одобрять"
//	@Failure		500				{object}	DefaultResponse[error]	"Внутренняя ошибка сервера"
//	@Router			/reservation/approve [patch]
func (h *httpDelivery) approveReservation(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.approveReservation")
	defer span.End()

	var req approveReservationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	role, ok := c.Get("role").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("failed to get role from context"))
	}

	if protopb.Role_ADMIN.String() != role || protopb.Role_ADMIN.String() != role {
		return c.JSON(http.StatusForbidden, ErrorResponse("only admins can approve reservations"))
	}

	if err := h.service.Reservation.ApproveReservation(ctx, req.ReservationId); err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusNoContent, DefaultResponse[string]{})
}

// createReservation godoc
//
//	@Summary		Create reservation
//	@Description	Создаёт бронирование на указанный период для текущего пользователя.
//	@Tags			reservation
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		createReservationRequest	true	"Create reservation request"
//	@Success		201		{object}	DefaultResponse[string]		"Успех"
//	@Failure		400		{object}	DefaultResponse[error]		"Невалидный запрос"
//	@Failure		401		{object}	DefaultResponse[error]		"Неавторизован"
//	@Failure		403		{object}	DefaultResponse[error]		"Нет доступа"
//	@Failure		500		{object}	DefaultResponse[error]		"Внутренняя ошибка сервера"
//	@Router			/reservation [post]
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

	userId, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("failed to get user_id from context"))
	}

	parsed, err := uuid.Parse(userId)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("invalid user id"))
	}

	username, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("failed to get username from context"))
	}

	role, ok := c.Get("role").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("failed to get role from context"))
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrorResponse("invalid date format (YYYY-MM-DD)"))
	}

	positions := make([]int16, 0, len(req.TimeSlots))
	for _, p := range req.TimeSlots {
		positions = append(positions, int16(p))
	}

	if err = h.service.Reservation.MakeReservation(ctx, parsed, parsedDate, positions, req.PeopleNum, role, username); err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusCreated, DefaultResponse[string]{
		Status: "success",
		Data:   "",
	})
}

// getFreeSlots godoc
//
//	@Summary		Get free time slots
//	@Description	Возвращает список свободных временных интервалов в заданном диапазоне.
//	@Description	Свободный слот — это интервал времени, не пересекающийся ни с одной резервацией.
//	@Tags			reservation
//	@Security		BearerAuth
//	@Produce		json
//	@Param			from	query		string									true	"Start datetime (RFC3339)"	example(2026-01-07T00:00:00Z)
//	@Param			to		query		string									true	"End datetime (RFC3339)"	example(2026-01-07T23:59:00Z)
//	@Success		200		{object}	DefaultResponse[getFreeSlotsResponse]	"List of free time intervals"
//	@Failure		400		{object}	DefaultResponse[error]					"Invalid request"
//	@Failure		401		{object}	DefaultResponse[error]					"Unauthorized"
//	@Failure		500		{object}	DefaultResponse[error]					"Internal server error"
//	@Router			/reservation/free-slots [get]
func (h *httpDelivery) getFreeSlots(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.getFreeSlots")
	defer span.End()

	var req getReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err := validate.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if req.Date == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse("date is required"))
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid date format (YYYY-MM-DD)"))
	}

	free, pairs, err := h.service.Reservation.GetFreeSlots(ctx, parsedDate)
	if err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusOK, DefaultResponse[getFreeSlotsResponse]{
		Data:   getFreeSlotsResponse{Free: free, Pairs: pairs},
		Status: "success",
	})
}

type getReservationRequest struct {
	Date string `query:"date" validate:"required"`
}

type getFreeSlotsResponse struct {
	Free  []model.DailySlot `json:"free_slots"`
	Pairs [][2]int16        `json:"free_pairs"`
}

type approveReservationRequest struct {
	ReservationId uuid.UUID `query:"reservation_id" validate:"required,uuid4"`
}

type createReservationRequest struct {
	TimeSlots []int  `json:"time_slots" validate:"required"`
	PeopleNum uint8  `json:"people_num" validate:"required,gt=1,lt=12"`
	Date      string `json:"date" validate:"required"` //2026-01-17
}
