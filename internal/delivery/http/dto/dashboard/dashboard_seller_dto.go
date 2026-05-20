package dashboard

import "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/others"

type SellerDashboardResponse struct {
	Seller      SellerInfo       `json:"seller"`
	Stats       SellerStats      `json:"stats"`
	Orders      SellerOrdersPage `json:"orders"`
	TopProducts []ProductSummary `json:"top_products"`
}

type SellerInfo struct {
	SellerID    string `json:"seller_id"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
}

type SellerStats struct {
	TotalRevenueBTC float64 `json:"total_revenue_btc"`
	TotalOrders     int64   `json:"total_orders"`
	PendingOrders   int64   `json:"pending_orders"`
	PaidOrders      int64   `json:"paid_orders"`
	ShippedOrders   int64   `json:"shipped_orders"`
	DeliveredOrders int64   `json:"delivered_orders"`
	CancelledOrders int64   `json:"cancelled_orders"`
	TotalProducts   int64   `json:"total_products"`
	AverageRating   float64 `json:"average_rating"`
	TotalReviews    int64   `json:"total_reviews"`
}

type SellerOrdersPage struct {
	Pagination others.Pagination `json:"pagination"`
	Items      []SellerOrder     `json:"items"`
}

type SellerOrder struct {
	OrderID       string          `json:"order_id"`
	BuyerID       string          `json:"buyer_id"`
	BuyerName     string          `json:"buyer_name"`
	BuyerEmail    string          `json:"buyer_email"`
	Status        string          `json:"status"`
	TotalBTC      float64         `json:"total_btc"`
	PaymentStatus string          `json:"payment_status"`
	Items         []OrderItemInfo `json:"items"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}

type OrderItemInfo struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	PriceBTC    float64 `json:"price_btc"`
	SubtotalBTC float64 `json:"subtotal_btc"`
}

type ProductSummary struct {
	ProductID  string  `json:"product_id"`
	Name       string  `json:"name"`
	TotalSold  int64   `json:"total_sold"`
	RevenueBTC float64 `json:"revenue_btc"`
	Stock      int     `json:"stock"`
	AvgRating  float64 `json:"avg_rating"`
}
