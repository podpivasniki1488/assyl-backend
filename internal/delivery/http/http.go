package http

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/podpivasniki1488/assyl-backend/docs"
	"github.com/podpivasniki1488/assyl-backend/internal/service"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type httpDelivery struct {
	echoApp   *echo.Echo
	logger    *slog.Logger
	service   *service.Service
	tracer    trace.Tracer
	jwtSecret string
}

func NewHTTPDelivery(logger *slog.Logger, s *service.Service, jwtSecret string) Http {
	return &httpDelivery{
		echoApp:   echo.New(),
		logger:    logger,
		service:   s,
		tracer:    otel.Tracer("httpDelivery"),
		jwtSecret: jwtSecret,
	}
}

var validate = validator.New()

func (h *httpDelivery) Start(port string) {
	h.registerHandler()

	h.echoApp.Use(middleware.CORS())

	h.echoApp.Debug = true

	if err := h.echoApp.Start(port); err != nil {
		h.logger.Error("Failed to start server", "error", err)
	}
}

func (h *httpDelivery) registerHandler() {
	v1 := h.echoApp.Group("/v1")

	v1.GET("/swagger/*", echoSwagger.WrapHandler)

	h.registerAuthHandlers(v1)
	h.registerApartmentHandlers(v1)
	h.registerReservationHandlers(v1)
	h.registerChannelHandlers(v1)
	h.registerFeedbackHandlers(v1)
}

func (h *httpDelivery) Stop(ctx context.Context) {
	if err := h.echoApp.Shutdown(ctx); err != nil {
		h.logger.Error("Failed to shutdown server", "error", err)
	}

	h.logger.Info("Shutting down server")
}
