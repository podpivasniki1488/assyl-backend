package http

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/protopb"
)

func (h *httpDelivery) registerChannelHandlers(v1 *echo.Group) {
	channel := v1.Group("/channel")

	channel.Use(h.registerJWTMiddleware())
	channel.GET("", h.getChannelMessages)
	channel.POST("", h.sendMessage, h.getJWTData())
}

// getChannelMessages godoc
//
//	@Summary		Get channel messages
//	@Description	Returns channel messages within the time period.
//	@Tags			Channel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			from	query		string	true	"From datetime (RFC3339)"	example(2025-12-01T00:00:00Z)
//	@Param			to		query		string	true	"To datetime (RFC3339)"		example(2025-12-02T00:00:00Z)
//	@Success		200		{object}	DefaultResponse[[]model.ChannelMessage]
//	@Failure		400		{object}	DefaultResponse[error]
//	@Failure		401		{object}	DefaultResponse[error]
//	@Failure		500		{object}	DefaultResponse[error]
//	@Router			/channel [get]
func (h *httpDelivery) getChannelMessages(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.getChannel")
	defer span.End()

	var req getChannelMessages
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err := validate.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	res, err := h.service.Channel.GetByTimePeriod(ctx, req.From, req.To)
	if err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusOK, DefaultResponse[[]model.ChannelMessage]{
		Status: "success",
		Data:   res,
	})
}

// sendMessage godoc
//
//	@Summary		Send channel message
//	@Description	Sends a message to the channel. Only ADMIN or GOD can send.
//	@Tags			Channel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			payload	body		sendChannelMessage	true	"Message payload"
//	@Success		204		{object}	DefaultResponse[string]
//	@Failure		400		{object}	DefaultResponse[error]
//	@Failure		401		{object}	DefaultResponse[error]
//	@Failure		403		{object}	DefaultResponse[error]
//	@Failure		500		{object}	DefaultResponse[error]
//	@Router			/channel [post]
func (h *httpDelivery) sendMessage(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.sendMessage")
	defer span.End()

	var req sendChannelMessage
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err := validate.Struct(&req); err != nil {
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

	role, ok := c.Get("role").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("failed to get role from context"))
	}

	if protopb.Role_ADMIN.String() != role && protopb.Role_GOD.String() != role {
		return c.JSON(http.StatusForbidden, ErrorResponse("only admins and gods can send messages"))
	}

	if err = h.service.Channel.SendChannelMessage(ctx, model.ChannelMessage{
		AuthorId:  parsed,
		Text:      req.Message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusNoContent, DefaultResponse[string]{})
}

type sendChannelMessage struct {
	Message string `json:"message" validate:"required"`
}

type getChannelMessages struct {
	From time.Time `query:"from" validate:"required"`
	To   time.Time `query:"to" validate:"required,gtfield=From"`
}
