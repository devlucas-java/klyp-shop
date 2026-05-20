package socket

import (
	"sync"

	"github.com/devlucas-java/klyp-shop/internal/infrastructure/observability/metrics"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type Client struct {
	UserID id.UUID
	Send   chan []byte
}

type Hub struct {
	mu      sync.RWMutex
	clients map[id.UUID]*Client
	log     *logger.Logger
	metric  *metrics.Metric
}

func NewHub(log *logger.Logger, metric *metrics.Metric) *Hub {
	return &Hub{
		clients: make(map[id.UUID]*Client),
		log:     log,
		metric:  metric,
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c.UserID] = c
	h.metric.WebSocketConnections.Inc()
	h.log.Infof("Hub: user %s connected", c.UserID)
}

func (h *Hub) Unregister(userID id.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if c, ok := h.clients[userID]; ok {
		close(c.Send)
		delete(h.clients, userID)
		h.metric.WebSocketConnections.Dec()
		h.log.Infof("Hub: user %s disconnected", userID)
	}
}

func (h *Hub) Send(receiverID id.UUID, msg []byte) {
	h.mu.RLock()
	c, ok := h.clients[receiverID]
	h.mu.RUnlock()
	if ok {
		select {
		case c.Send <- msg:
		default:
			h.log.Warnf("Hub: send buffer full for user %s, dropping message", receiverID)
		}
	}
}

func (h *Hub) IsOnline(userID id.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}
