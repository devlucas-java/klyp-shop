package module

import (
	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/router"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/btcpay"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
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
) chi.Router {
	paymentRepository := database.NewBitcoinPaymentDB(db, log)
	orderRepository := database.NewOrderDB(db, log)
	userRepository := database.NewUserDB(db, log)

	btcpayClient := btcpay.NewClient(btcpayBaseURL, btcpayStoreID, btcpayAPIKey)

	paymentService := service.NewPaymentService(log, paymentRepository, orderRepository, btcpayClient, webhookSecret)
	paymentHandler := handler.NewPaymentHandler(paymentService, log)
	paymentRouter := router.NewPaymentRouter(jwtService, paymentHandler, log, userRepository)

	r := chi.NewRouter()
	paymentRouter.RegisterPaymentRoutes(r)
	return r
}
