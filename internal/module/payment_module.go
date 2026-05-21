package module

import (
	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/adapter"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/router"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/client/btcpay"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/observability/metrics"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func InitPaymentModule(
	db *gorm.DB,
	log *logger.Logger,
	jwtService *jwt.JWTService,
	btcpayBaseURL, btcpayStoreID, btcpayAPIKey, webhookSecret string,
	metric *metrics.Metric,
) chi.Router {
	paymentRepository := database.NewBitcoinPaymentDB(db)
	orderRepository := database.NewOrderDB(db)
	userRepository := database.NewUserDB(db)

	gateway := btcpay.NewBTCPayAdapter(btcpayBaseURL, btcpayStoreID, btcpayAPIKey, webhookSecret)

	paymentService := service.NewPaymentService(log, paymentRepository, orderRepository, gateway, metric)
	paymentHandler := handler.NewPaymentHandler(paymentService, log)
	adapter := adapter.NewAdapter(log)
	paymentRouter := router.NewPaymentRouter(jwtService, paymentHandler, log, userRepository, adapter)

	r := chi.NewRouter()
	paymentRouter.RegisterPaymentRoutes(r)
	return r
}
