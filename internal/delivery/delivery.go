package delivery

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/podpivasniki1488/assyl-backend/internal/delivery/http"
	"github.com/podpivasniki1488/assyl-backend/internal/service"
)

type Delivery struct {
	Http http.Http
}

func NewDelivery(logger *slog.Logger, e *echo.Echo, service *service.Service, jwtSecret string) *Delivery {
	return &Delivery{
		Http: http.NewHTTPDelivery(logger, e, service, jwtSecret),
	}
}
