package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

func (h *httpDelivery) registerAuthHandlers(v1 *echo.Group) {
	auth := v1.Group("/auth")

	auth.POST("/register", h.register)
	auth.POST("/confirm", h.confirm)
	auth.POST("/login", h.login)
}

// confirm godoc
// @Summary      Confirm registration
// @Description  Подтверждает пользователя по OTP-коду.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      confirmRequest       true  "Confirm request"
// @Success      204      "Подтверждение прошло успешно"
// @Failure      400      {object}  DefaultResponse[error]  "Невалидный запрос"
// @Failure      500      {object}  DefaultResponse[error]  "Внутренняя ошибка сервера"
// @Router       /auth/confirm [post]
func (h *httpDelivery) confirm(c echo.Context) error {
	ctx := c.Request().Context()

	ctx, span := h.tracer.Start(ctx, "httpDelivery.confirm")
	defer span.End()

	var confirm confirmRequest
	if err := c.Bind(&confirm); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err := validator.New().Struct(&confirm); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err := h.service.Auth.Confirm(ctx, confirm.Username, confirm.OtpCode); err != nil {
		return HandleErrResponse(c, err)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *httpDelivery) login(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.login")
	defer span.End()

	var login loginRequest
	if err := c.Bind(&login); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err := validator.New().Struct(&login); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	token, err := h.service.Auth.Login(ctx, model.User{
		Username: login.Username,
		Password: login.Password,
	})
	if err != nil {
		return HandleErrResponse(c, err)
	}

	return c.JSON(http.StatusOK, DefaultResponse[string]{
		Status: "ok",
		Data:   token,
	})
}

// register godoc
// @Summary      Register new user
// @Description  Регистрирует нового пользователя и отправляет OTP на указанный username (телефон/почта).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      registerRequest      true  "Register request"
// @Success      200      {object}  DefaultResponse[string] "Успешная регистрация"
// @Failure      400      {object}  DefaultResponse[error]  "Невалидный запрос"
// @Failure      500      {object}  DefaultResponse[error]  "Внутренняя ошибка сервера"
// @Router       /auth/register [post]
func (h *httpDelivery) register(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.register")
	defer span.End()

	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err := validator.New().Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err := h.service.Auth.Register(ctx,
		model.User{
			Username: req.Username,
			Password: req.Password,
		}); err != nil {
		return HandleErrResponse(c, err)
	}

	return c.JSON(http.StatusOK, DefaultResponse[string]{
		Status: "ok",
	})
}

func (h *httpDelivery) logout(c echo.Context) error {
	return nil
}

type confirmRequest struct {
	Username string `json:"username" validate:"required"`
	OtpCode  string `json:"otp_code" validate:"required"`
}

type registerRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=32"`
	Password  string `json:"password" validate:"required,min=8,max=32"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type loginRequest struct {
	Username string `json:"Username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
