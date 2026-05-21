package dashboard

import (
	"github.com/devlucas-java/klyp-shop/pkg/pagination"
)

type AdminDashboardResponse struct {
	Stats      AdminStats      `json:"stats"`
	Orders     AdminOrdersPage `json:"orders"`
	TopSellers []SellerRanking `json:"top_sellers"`
}

type AdminStats struct {
	TotalRevenueBTC float64        `json:"total_revenue_btc"`
	TotalOrders     int64          `json:"total_orders"`
	TotalUsers      int64          `json:"total_users"`
	TotalSellers    int64          `json:"total_sellers"`
	TotalProducts   int64          `json:"total_products"`
	OrdersByStatus  OrdersByStatus `json:"orders_by_status"`
}

type OrdersByStatus struct {
	Pending   int64 `json:"pending"`
	Paid      int64 `json:"paid"`
	Shipped   int64 `json:"shipped"`
	Delivered int64 `json:"delivered"`
	Cancelled int64 `json:"cancelled"`
}

type AdminOrdersPage struct {
	Pagination pagination.OutPutPagination `json:"pagination"`
	Items      []AdminOrder                `json:"items"`
}

type AdminOrder struct {
	OrderID       string           `json:"order_id"`
	BuyerID       string           `json:"buyer_id"`
	BuyerName     string           `json:"buyer_name"`
	BuyerEmail    string           `json:"buyer_email"`
	Status        string           `json:"status"`
	TotalBTC      float64          `json:"total_btc"`
	PaymentStatus string           `json:"payment_status"`
	Items         []AdminOrderItem `json:"items"`
	CreatedAt     string           `json:"created_at"`
	UpdatedAt     string           `json:"updated_at"`
}

type AdminOrderItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	SellerID    string  `json:"seller_id"`
	SellerName  string  `json:"seller_name"`
	Quantity    int     `json:"quantity"`
	PriceBTC    float64 `json:"price_btc"`
	SubtotalBTC float64 `json:"subtotal_btc"`
}

type SellerRanking struct {
	SellerID    string  `json:"seller_id"`
	DisplayName string  `json:"display_name"`
	TotalOrders int64   `json:"total_orders"`
	RevenueBTC  float64 `json:"revenue_btc"`
	TotalSold   int64   `json:"total_sold"`
}
