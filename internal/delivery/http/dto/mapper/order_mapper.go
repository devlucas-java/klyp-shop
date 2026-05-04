package mapper

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dorder"
	dorderitem "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dorder_item"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type OrderMapper struct{}

func NewOrderMapper() *OrderMapper {
	return &OrderMapper{}
}

func (m *OrderMapper) OrderToResponse(order *entity.Order) *dorder.OrderResponse {
	if order == nil {
		return nil
	}

	items := make([]dorderitem.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = dorderitem.OrderItemResponse{
			ID:        item.ID.String(),
			ProductID: item.ProductID.String(),
			Quantity:  item.Quantity,
			PriceBTC:  item.PriceBTC,
			Subtotal:  item.Subtotal(),
		}
	}

	return &dorder.OrderResponse{
		ID:        order.ID.String(),
		UserID:    order.UserID.String(),
		AddressID: order.AddressID.String(),
		Status:    string(order.Status),
		TotalBTC:  order.TotalBTC,
		Items:     items,
		CreatedAt: order.CreatedAt.String(),
		UpdatedAt: order.UpdatedAt.String(),
	}
}

func (m *OrderMapper) OrdersToResponses(orders []*entity.Order) []*dorder.OrderResponse {
	if orders == nil || len(orders) == 0 {
		return []*dorder.OrderResponse{}
	}

	responses := make([]*dorder.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = m.OrderToResponse(order)
	}
	return responses
}
