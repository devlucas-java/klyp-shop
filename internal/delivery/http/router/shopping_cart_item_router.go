package router

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/adapter"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/go-chi/chi"
)

type ShoppingCartItemRouter struct {
	shoppingCartItemHandler *handler.ShoppingCartItemHandler
	adapter                 *adapter.Adapter
}

func NewShoppingCartItemRouter(
	sh *handler.ShoppingCartItemHandler,
	adapter *adapter.Adapter,
) *ShoppingCartItemRouter {
	return &ShoppingCartItemRouter{
		shoppingCartItemHandler: sh,
		adapter:                 adapter,
	}
}

func (s *ShoppingCartItemRouter) RegisterShoppingCartItemRoutes(r chi.Router) {
	r.Post("/items", s.adapter.Adapt(s.shoppingCartItemHandler.AddItem))
	r.Patch("/items/{id}", s.adapter.Adapt(s.shoppingCartItemHandler.UpdateItem))
	r.Delete("/items/{id}", s.adapter.Adapt(s.shoppingCartItemHandler.RemoveItem))
}
