package http

import (
	"errors"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

type DefaultResponse[T any] struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
	Data         T      `json:"data"`
}

func ErrorResponse(errMsg string) DefaultResponse[error] {
	return DefaultResponse[error]{
		Status:       "error",
		ErrorMessage: errMsg,
		Data:         nil,
	}
}

func (h *httpDelivery) handleErrResponse(c echo.Context, err error) error {
	var appErr model.AppError
	if errors.As(err, &appErr) {
		return c.JSON(appErr.HttpStatusCode, ErrorResponse(appErr.Error()))
	}

	h.logger.Error("internal error %s", err.Error())

	return c.JSON(http.StatusInternalServerError, ErrorResponse("internal error"))
}

func (h *httpDelivery) registerJWTMiddleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(h.jwtSecret),
	})
}
