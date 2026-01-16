package pkg

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type HubMessage struct {
	ChatID uuid.UUID `json:"chat_id"`
	Text   string    `json:"text"`
}

type Hub struct {
	mux  sync.RWMutex
	subs map[uuid.UUID]chan HubMessage //userId hubMessage
}

var (
	ErrSubNotFound = errors.New("channel not found")
)

func NewHub() *Hub {
	return &Hub{subs: make(map[uuid.UUID]chan HubMessage)}
}

func (h *Hub) Subscribe(userId uuid.UUID) (chan HubMessage, func()) {
	h.mux.Lock()
	defer h.mux.Unlock()

	ch := make(chan HubMessage)
	h.subs[userId] = ch

	unsubscribe := func() {
		h.mux.Lock()
		delete(h.subs, userId)
		close(ch)
		h.mux.Unlock()
	}

	return ch, unsubscribe
}

func (h *Hub) Publish(userId uuid.UUID, msg HubMessage) error {
	h.mux.RLock()
	defer h.mux.RUnlock()

	ch, ok := h.subs[userId]
	if !ok {
		return ErrSubNotFound
	}

	ch <- msg

	return nil
}
