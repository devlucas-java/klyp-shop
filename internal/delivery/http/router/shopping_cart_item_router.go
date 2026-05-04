package router

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/adapter"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/go-chi/chi"
)

type ShoppingCartItemRouter struct {
	shoppingCartItemHandler *handler.ShoppingCartItemHandler
}

func NewShoppingCartItemRouter(shoppingCartItemHandler *handler.ShoppingCartItemHandler) *ShoppingCartItemRouter {
	return &ShoppingCartItemRouter{
		shoppingCartItemHandler: shoppingCartItemHandler,
	}
}

func (r *ShoppingCartItemRouter) RegisterShoppingCartItemRoutes(protect chi.Router) {
	protect.Post("/items", adapter.Adapt(r.shoppingCartItemHandler.AddItem))
	protect.Patch("/items/{id}", adapter.Adapt(r.shoppingCartItemHandler.UpdateItem))
	protect.Delete("/items/{id}", adapter.Adapt(r.shoppingCartItemHandler.RemoveItem))
}
