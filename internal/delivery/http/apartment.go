package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

func (h *httpDelivery) registerApartmentHandlers(v1 *echo.Group) {
	apartment := v1.Group("/apartment")
	apartment.Use(h.registerJWTMiddleware())

	apartment.POST("/create", h.createApartment)
	apartment.GET("", h.getApartment)
}

// getApartment godoc
//
//	@Summary		Get apartment
//	@Description	Возвращает квартиру по этажу и номеру двери.
//	@Tags			apartment
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			floor		query		int									true	"Этаж"
//	@Param			door_number	query		int									true	"Номер двери"
//	@Success		200			{object}	DefaultResponse[model.Apartment]	"Квартира найдена"
//	@Failure		400			{object}	DefaultResponse[error]				"Невалидный запрос"
//	@Failure		401			{object}	DefaultResponse[error]				"Неавторизован"
//	@Failure		404			{object}	DefaultResponse[error]				"Квартира не найдена"
//	@Failure		500			{object}	DefaultResponse[error]				"Внутренняя ошибка сервера"
//	@Router			/apartment [get]
func (h *httpDelivery) getApartment(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.getApartment")
	defer span.End()

	var req getApartmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if req.DoorNumber == nil || req.Floor == nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("Missing mandatory field: doorNumber, floor"))
	}

	res, err := h.service.Apartment.GetApartment(ctx, model.Apartment{
		DoorNumber: *req.DoorNumber,
		Floor:      *req.Floor,
	})
	if err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusOK, DefaultResponse[model.Apartment]{
		Status: "ok",
		Data:   *res,
	})
}

// createApartment godoc
//
//	@Summary		Create apartment
//	@Description	Создаёт новую квартиру.
//	@Tags			apartment
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			request	body		createApartmentRequest	true	"Create apartment request"
//	@Success		201		{object}	DefaultResponse[string]	"Квартира успешно создана"
//	@Failure		400		{object}	DefaultResponse[error]	"Невалидный запрос"
//	@Failure		401		{object}	DefaultResponse[error]	"Неавторизован"
//	@Failure		409		{object}	DefaultResponse[error]	"Квартира уже существует (конфликт)"
//	@Failure		500		{object}	DefaultResponse[error]	"Внутренняя ошибка сервера"
//	@Router			/apartment/create [post]
func (h *httpDelivery) createApartment(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.createApartment")
	defer span.End()

	var req createApartmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err := validate.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err := h.service.Apartment.CreateApartment(ctx, model.Apartment{
		Floor:      req.Floor,
		DoorNumber: req.DoorNum,
	}); err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusCreated, DefaultResponse[string]{})
}

type createApartmentRequest struct {
	Floor   uint8  `json:"floor" validate:"required"`
	DoorNum uint16 `json:"door_num" validate:"required"`
}

type getApartmentRequest struct {
	Floor      *uint8  `query:"floor"`
	DoorNumber *uint16 `query:"door_number"`
}
