package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/protopb"
)

func (h *httpDelivery) registerOrderHandlers(v1 *echo.Group) {
	order := v1.Group("/order")

	order.Use(h.registerJWTMiddleware())
	order.GET("", h.getOrders, h.getJWTData())
	order.POST("", h.createOrder, h.getJWTData())
}

func (h *httpDelivery) createOrder(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.createOrder")
	defer span.End()

	var req createOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	userId, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("failed to get user_id from context"))
	}

	parsed, err := uuid.Parse(userId)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("invalid user id"))
	}

	val, ok := protopb.OrderType_value[req.OrderType]
	if !ok {
		return c.JSON(http.StatusUnprocessableEntity, ErrorResponse("invalid order type"))
	}

	orderType := protopb.OrderType(val)

	if err = h.service.Order.OrderService(ctx, &model.Order{
		UserID:    parsed,
		Text:      req.Text,
		OrderType: orderType,
	}); err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusNoContent, DefaultResponse[string]{})
}

type createOrderRequest struct {
	OrderType string `json:"order_type"`
	Text      string `json:"text"`
}

func (h *httpDelivery) getOrders(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.getOrders")
	defer span.End()

	var req getOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err := validate.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	role, ok := c.Get("role").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("failed to get role from context"))
	}

	userId, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("failed to get user_id from context"))
	}

	parsed, err := uuid.Parse(userId)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("invalid user id"))
	}

	resp, err := h.service.Order.GetUserOrders(ctx, &model.GetOrderRequest{
		UserID:    parsed,
		OrderType: req.OrderType,
		Text:      req.Text,
	}, role)
	if err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusOK, DefaultResponse[[]model.Order]{
		Data:   resp,
		Status: "success",
	})
}

type getOrderRequest struct {
	ID        uuid.UUID `query:"id"`
	UserID    uuid.UUID `query:"user_id"`
	OrderType string    `query:"order_type"`
	Text      string    `query:"text"`
}
