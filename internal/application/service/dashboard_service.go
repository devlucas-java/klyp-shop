package service

import (
	"math"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/ddashboard"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type DashboardService struct {
	log                 *logger.Logger
	userRepository      repository.UserRepository
	orderRepository     repository.OrderRepository
	dashboardRepository repository.DashboardRepository
}

func NewDashboardService(
	log *logger.Logger,
	userRepository repository.UserRepository,
	orderRepository repository.OrderRepository,
	dashboardRepository repository.DashboardRepository,
) *DashboardService {
	return &DashboardService{
		log:                 log,
		userRepository:      userRepository,
		orderRepository:     orderRepository,
		dashboardRepository: dashboardRepository,
	}
}

func (s *DashboardService) GetSellerDashboard(auth *entity.User, page, size int, statusFilter string) (*ddashboard.SellerDashboardResponse, error) {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return nil, errors.ErrNotFound("User", err)
	}
	if err := user.EnsureSeller(); err != nil {
		return nil, err
	}

	sellerID := user.Seller.ID

	statusCounts, err := s.dashboardRepository.CountOrdersByStatusForSeller(sellerID)
	if err != nil {
		return nil, errors.ErrDatabase("failed to count orders", err)
	}

	revenue, err := s.dashboardRepository.SumRevenueForSeller(sellerID)
	if err != nil {
		return nil, errors.ErrDatabase("failed to sum revenue", err)
	}

	productCount, err := s.dashboardRepository.CountProductsForSeller(sellerID)
	if err != nil {
		return nil, errors.ErrDatabase("failed to count products", err)
	}

	avgRating, totalReviews, err := s.dashboardRepository.AvgRatingForSeller(sellerID)
	if err != nil {
		return nil, errors.ErrDatabase("failed to get rating", err)
	}

	topProductRows, err := s.dashboardRepository.TopProductsForSeller(sellerID, 10)
	if err != nil {
		return nil, errors.ErrDatabase("failed to get top products", err)
	}

	orders, total, err := s.orderRepository.FindBySellerIDPaginated(sellerID, page, size, statusFilter)
	if err != nil {
		return nil, errors.ErrDatabase("failed to fetch orders", err)
	}

	stats := buildSellerStats(statusCounts, revenue, productCount, avgRating, totalReviews)

	topProducts := make([]ddashboard.ProductSummary, len(topProductRows))
	for i, r := range topProductRows {
		topProducts[i] = ddashboard.ProductSummary{
			ProductID:  r.ProductID,
			Name:       r.Name,
			TotalSold:  r.TotalSold,
			RevenueBTC: r.RevenueBTC,
			Stock:      r.Stock,
		}
	}

	orderItems := buildSellerOrders(orders)

	return &ddashboard.SellerDashboardResponse{
		Seller: ddashboard.SellerInfo{
			SellerID:    sellerID.String(),
			DisplayName: user.Seller.DisplayName,
			Bio:         user.Seller.Bio,
		},
		Stats: stats,
		Orders: ddashboard.SellerOrdersPage{
			Pagination: paginate(page, size, total),
			Items:      orderItems,
		},
		TopProducts: topProducts,
	}, nil
}

func (s *DashboardService) GetAdminDashboard(page, size int, statusFilter string) (*ddashboard.AdminDashboardResponse, error) {
	totalUsers, err := s.dashboardRepository.CountAllUsers()
	if err != nil {
		return nil, errors.ErrDatabase("failed to count users", err)
	}

	totalSellers, err := s.dashboardRepository.CountAllSellers()
	if err != nil {
		return nil, errors.ErrDatabase("failed to count sellers", err)
	}

	totalProducts, err := s.dashboardRepository.CountAllProducts()
	if err != nil {
		return nil, errors.ErrDatabase("failed to count products", err)
	}

	statusCounts, err := s.dashboardRepository.CountAllOrdersByStatus()
	if err != nil {
		return nil, errors.ErrDatabase("failed to count orders by status", err)
	}

	totalRevenue, err := s.dashboardRepository.SumTotalRevenue()
	if err != nil {
		return nil, errors.ErrDatabase("failed to sum revenue", err)
	}

	topSellerRows, err := s.dashboardRepository.TopSellersByRevenue(10)
	if err != nil {
		return nil, errors.ErrDatabase("failed to get top sellers", err)
	}

	orders, total, err := s.orderRepository.FindAllPaginated(page, size, statusFilter)
	if err != nil {
		return nil, errors.ErrDatabase("failed to fetch orders", err)
	}

	ordersByStatus := buildOrdersByStatus(statusCounts)
	var totalOrders int64
	totalOrders = ordersByStatus.Pending + ordersByStatus.Paid + ordersByStatus.Shipped +
		ordersByStatus.Delivered + ordersByStatus.Cancelled

	topSellers := make([]ddashboard.SellerRanking, len(topSellerRows))
	for i, r := range topSellerRows {
		topSellers[i] = ddashboard.SellerRanking{
			SellerID:    r.SellerID,
			DisplayName: r.DisplayName,
			TotalOrders: r.TotalOrders,
			RevenueBTC:  r.RevenueBTC,
			TotalSold:   r.TotalSold,
		}
	}

	adminOrders := buildAdminOrders(orders)

	return &ddashboard.AdminDashboardResponse{
		Stats: ddashboard.AdminStats{
			TotalRevenueBTC: totalRevenue,
			TotalOrders:     totalOrders,
			TotalUsers:      totalUsers,
			TotalSellers:    totalSellers,
			TotalProducts:   totalProducts,
			OrdersByStatus:  ordersByStatus,
		},
		Orders: ddashboard.AdminOrdersPage{
			Pagination: paginate(page, size, total),
			Items:      adminOrders,
		},
		TopSellers: topSellers,
	}, nil
}

func buildSellerStats(
	counts []repository.OrderStatusCount,
	revenue float64,
	products int64,
	avgRating float64,
	totalReviews int64,
) ddashboard.SellerStats {
	s := ddashboard.SellerStats{
		TotalRevenueBTC: revenue,
		TotalProducts:   products,
		AverageRating:   avgRating,
		TotalReviews:    totalReviews,
	}
	for _, c := range counts {
		s.TotalOrders += c.Count
		switch c.Status {
		case string(entity.OrderStatusPending):
			s.PendingOrders = c.Count
		case string(entity.OrderStatusPaid):
			s.PaidOrders = c.Count
		case string(entity.OrderStatusShipped):
			s.ShippedOrders = c.Count
		case string(entity.OrderStatusDelivered):
			s.DeliveredOrders = c.Count
		case string(entity.OrderStatusCancelled):
			s.CancelledOrders = c.Count
		}
	}
	return s
}

func buildOrdersByStatus(counts []repository.OrderStatusCount) ddashboard.OrdersByStatus {
	var s ddashboard.OrdersByStatus
	for _, c := range counts {
		switch c.Status {
		case string(entity.OrderStatusPending):
			s.Pending = c.Count
		case string(entity.OrderStatusPaid):
			s.Paid = c.Count
		case string(entity.OrderStatusShipped):
			s.Shipped = c.Count
		case string(entity.OrderStatusDelivered):
			s.Delivered = c.Count
		case string(entity.OrderStatusCancelled):
			s.Cancelled = c.Count
		}
	}
	return s
}

func buildSellerOrders(orders []*entity.Order) []ddashboard.SellerOrder {
	result := make([]ddashboard.SellerOrder, 0, len(orders))
	for _, o := range orders {
		paymentStatus := "none"
		if o.BitcoinPayment != nil {
			paymentStatus = string(o.BitcoinPayment.Status)
		}

		items := make([]ddashboard.OrderItemInfo, 0, len(o.Items))
		for _, item := range o.Items {
			items = append(items, ddashboard.OrderItemInfo{
				ProductID:   item.ProductID.String(),
				ProductName: item.Product.Name,
				Quantity:    item.Quantity,
				PriceBTC:    item.PriceBTC,
				SubtotalBTC: item.PriceBTC * float64(item.Quantity),
			})
		}

		result = append(result, ddashboard.SellerOrder{
			OrderID:       o.ID.String(),
			BuyerID:       o.UserID.String(),
			BuyerName:     o.User.Name,
			BuyerEmail:    o.User.Email,
			Status:        string(o.Status),
			TotalBTC:      o.TotalBTC,
			PaymentStatus: paymentStatus,
			Items:         items,
			CreatedAt:     o.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     o.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return result
}

func buildAdminOrders(orders []*entity.Order) []ddashboard.AdminOrder {
	result := make([]ddashboard.AdminOrder, 0, len(orders))
	for _, o := range orders {
		paymentStatus := "none"
		if o.BitcoinPayment != nil {
			paymentStatus = string(o.BitcoinPayment.Status)
		}

		items := make([]ddashboard.AdminOrderItem, 0, len(o.Items))
		for _, item := range o.Items {
			sellerID := item.Product.SellerID.String()
			sellerName := ""
			if item.Product.SellerID != (item.Product.SellerID) {
				sellerName = sellerID
			}
			items = append(items, ddashboard.AdminOrderItem{
				ProductID:   item.ProductID.String(),
				ProductName: item.Product.Name,
				SellerID:    sellerID,
				SellerName:  sellerName,
				Quantity:    item.Quantity,
				PriceBTC:    item.PriceBTC,
				SubtotalBTC: item.PriceBTC * float64(item.Quantity),
			})
		}

		result = append(result, ddashboard.AdminOrder{
			OrderID:       o.ID.String(),
			BuyerID:       o.UserID.String(),
			BuyerName:     o.User.Name,
			BuyerEmail:    o.User.Email,
			Status:        string(o.Status),
			TotalBTC:      o.TotalBTC,
			PaymentStatus: paymentStatus,
			Items:         items,
			CreatedAt:     o.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     o.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return result
}

func paginate(page, size int, total int64) ddashboard.Pagination {
	totalPages := int64(math.Ceil(float64(total) / float64(size)))
	if totalPages < 1 {
		totalPages = 1
	}
	return ddashboard.Pagination{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
	}
}
