package delivery

import (
	"log/slog"

	"github.com/podpivasniki1488/assyl-backend/internal/delivery/http"
	"github.com/podpivasniki1488/assyl-backend/internal/service"
	"go.opentelemetry.io/otel/trace"
)

type Delivery struct {
	Http http.Http
	// maybe ws
}

func NewDelivery(logger *slog.Logger, service *service.Service, trace trace.Tracer) *Delivery {
	return &Delivery{
		Http: http.NewHTTPDelivery(logger, service, trace),
	}
}
