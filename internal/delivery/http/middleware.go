package http

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (h *httpDelivery) getJWTData() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return c.JSON(http.StatusUnauthorized, "JWT token missing or invalid")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, "failed to cast claims")
			}

			username, ok := claims["username"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, "failed to cast username")
			}

			userId, ok := claims["user_id"].(uuid.UUID)
			if !ok {
				return c.JSON(http.StatusUnauthorized, "failed to cast user_id")
			}

			c.Set("username", username)
			c.Set("user_id", userId)

			return next(c)
		}
	}
}
