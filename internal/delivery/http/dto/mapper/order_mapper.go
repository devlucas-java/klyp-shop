package mapper

import (
	orderDTO "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/order"
	orderitemDTO "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/order_item"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type OrderMapper struct{}

func NewOrderMapper() *OrderMapper {
	return &OrderMapper{}
}

func (m *OrderMapper) OrderToResponse(order *entity.Order) *orderDTO.OrderResponse {
	if order == nil {
		return nil
	}

	items := make([]orderitemDTO.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = orderitemDTO.OrderItemResponse{
			ID:        item.ID.String(),
			ProductID: item.ProductID.String(),
			Quantity:  item.Quantity,
			PriceBTC:  item.PriceBTC,
			Subtotal:  item.Subtotal(),
		}
	}

	return &orderDTO.OrderResponse{
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

func (m *OrderMapper) OrdersToResponses(orders []*entity.Order) []*orderDTO.OrderResponse {
	if len(orders) == 0 {
		return []*orderDTO.OrderResponse{}
	}

	responses := make([]*orderDTO.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = m.OrderToResponse(order)
	}
	return responses
}
