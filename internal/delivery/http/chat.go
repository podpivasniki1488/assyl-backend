package http

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *httpDelivery) registerChatHandlers(v1 *echo.Group) {
	chat := v1.Group("/chat")
	chat.Use(h.registerJWTMiddleware())

	chat.POST("/send", h.sendMsgToChat, h.getJWTData())
	chat.POST("/start", h.startChat, h.getJWTData())

	ws := v1.Group("/ws")
	ws.Use(h.registerJWTMiddleware())

	ws.GET("/last-message", h.getUserLastMessages, h.getJWTData()) // listen for user-last messages (must be websocket)
}

func (h *httpDelivery) startChat(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.startChat")
	defer span.End()

	userId, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("user id not found in context"))
	}

	parsed, err := uuid.Parse(userId)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("invalid user id"))
	}

	var req startChatRequest
	if err = c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err = validate.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	res, err := h.service.Chat.StartNewChat(ctx, parsed, req.Members)
	if err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusOK, DefaultResponse[uuid.UUID]{
		Status: "ok",
		Data:   res,
	})
}

type startChatRequest struct {
	Members []uuid.UUID `json:"members" validate:"required"` // single member, means chat is private
}

func (h *httpDelivery) sendMsgToChat(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.sendMsgToChat")
	defer span.End()

	userId, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("role not found in context"))
	}

	parsed, err := uuid.Parse(userId)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("invalid user id"))
	}

	var req sendMsgToChatRequest
	if err = c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err = validate.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	if err = h.service.Chat.SendMessageToChat(ctx, req.Message, req.ChatId, parsed); err != nil {
		return h.handleErrResponse(c, err)
	}

	return c.JSON(http.StatusNoContent, DefaultResponse[string]{})
}

type sendMsgToChatRequest struct {
	ChatId        uuid.UUID  `json:"chat_id" validate:"required,uuid"`
	Message       string     `json:"message" validate:"required"`
	ReplyToUserId *uuid.UUID `json:"reply_to,omitempty"`
	// TODO: maybe add attachment links
}

// use when user wants to listen for last messages
func (h *httpDelivery) getUserLastMessages(c echo.Context) error {
	_, span := h.tracer.Start(c.Request().Context(), "httpDelivery.getUserLastMessages")
	defer span.End()

	userId, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("userId not found in context"))
	}

	parsed, err := uuid.Parse(userId)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("invalid user id"))
	}

	conn, err := Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		h.logger.Error("failed to upgrade websocket", "error", err)
		return err
	}
	defer conn.Close()

	subCh, unsub := h.service.Chat.SubscribeToUpdates(parsed)
	defer unsub()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err = conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-subCh:
			if !ok {
				return nil
			}

			if err = conn.WriteJSON(msg); err != nil {
				h.logger.Error("failed to write message", "error", err)
				return err
			}

		case <-done:
			return nil

		case <-ticker.C:
			if err = conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(5*time.Second)); err != nil {
				h.logger.Error("failed to send ping", "error", err)
				return nil
			}
		}
	}
}
