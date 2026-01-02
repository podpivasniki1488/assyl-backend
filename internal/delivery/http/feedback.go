package http

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/protopb"
)

func (h *httpDelivery) registerFeedbackHandlers(v1 *echo.Group) {
	feedback := v1.Group("/feedback")

	feedback.Use(h.registerJWTMiddleware())
	feedback.POST("", h.createFeedback, h.getJWTData())
	feedback.GET("", h.getFeedbacks, h.getJWTData())
}

// createFeedback godoc
//
//	@Summary		Create feedback
//	@Description	Creates feedback message from authorized user
//	@Tags			feedback
//	@Security		BearerAuth
//	@Param			payload	body		createFeedbackRequest	true	"Feedback payload"
//	@Success		200		{object}	DefaultResponse[string]
//	@Failure		400		{object}	DefaultResponse[error]
//	@Failure		401		{object}	DefaultResponse[error]
//	@Failure		500		{object}	DefaultResponse[error]
//	@Router			/feedback [post]
func (h *httpDelivery) createFeedback(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.createFeedback")
	defer span.End()

	var req createFeedbackRequest
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

	if err = h.service.Feedback.CreateFeedback(ctx, model.Feedback{
		UserId:       parsed,
		Text:         req.Message,
		FeedbackType: req.FeedbackType,
	}); err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusCreated, DefaultResponse[string]{})
}

type createFeedbackRequest struct {
	Message      string               `json:"message" validate:"required"`
	FeedbackType protopb.FeedbackType `json:"feedback_type" validate:"required"`
}

// getFeedbacks godoc
//
//	@Summary		Get feedbacks
//	@Description	Returns feedbacks list (only ADMIN or GOD)
//	@Tags			feedback
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			created_at_from	query		string	false	"Created from (RFC3339)"	format(date-time)	example(2026-01-01T00:00:00Z)
//	@Param			created_at_to	query		string	false	"Created to (RFC3339)"		format(date-time)	example(2026-01-02T00:00:00Z)
//	@Param			text			query		string	false	"Text filter"
//	@Param			feedback_type	query		string	false	"Feedback type"
//	@Success		200				{object}	DefaultResponse[string]
//	@Failure		400				{object}	DefaultResponse[error]
//	@Failure		401				{object}	DefaultResponse[error]
//	@Failure		500				{object}	DefaultResponse[error]
//	@Router			/feedback [get]
func (h *httpDelivery) getFeedbacks(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.getFeedbacks")
	defer span.End()

	var req getFeedbackRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	role, ok := c.Get("role").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("failed to get role from context"))
	}

	if protopb.Role_ADMIN.String() != role && protopb.Role_GOD.String() != role {
		return c.JSON(http.StatusForbidden, ErrorResponse("only admins and gods can send messages"))
	}

	userId, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("failed to get user_id from context"))
	}

	parsed, err := uuid.Parse(userId)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("invalid user id"))
	}

	txt := ""
	if req.Text != nil {
		txt = *req.Text
	}

	feedbackType := ""
	if req.FeedbackType != nil {
		feedbackType = *req.FeedbackType
	}

	res, err := h.service.Feedback.GetFeedbacks(ctx, model.GetFeedbackRequest{
		UserID:        parsed,
		CreatedAtFrom: req.CreatedAtFrom,
		CreatedAtTo:   req.CreatedAtTo,
		Text:          txt,
		FeedbackType:  feedbackType,
	})
	if err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusOK, DefaultResponse[[]model.Feedback]{
		Status: "ok",
		Data:   res,
	})
}

type getFeedbackRequest struct {
	CreatedAtFrom time.Time `query:"created_at_from"`
	CreatedAtTo   time.Time `query:"created_at_to"`
	Text          *string   `query:"text"`
	FeedbackType  *string   `query:"feedback_type"`
}
