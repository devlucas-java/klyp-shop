package router

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/adapter"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/go-chi/chi"
)

type OrderItemRouter struct {
	orderItemHandler *handler.OrderItemHandler
}

func NewOrderItemRouter(orderItemHandler *handler.OrderItemHandler) *OrderItemRouter {
	return &OrderItemRouter{
		orderItemHandler: orderItemHandler,
	}
}

func (r *OrderItemRouter) RegisterOrderItemRoutes(protect chi.Router) {
	protect.Get("/{id}/items", adapter.Adapt(r.orderItemHandler.GetOrderItems))
	protect.Get("/{id}/items/{itemId}", adapter.Adapt(r.orderItemHandler.GetOrderItem))
}
