package module

import (
	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/adapter"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/router"
	"github.com/devlucas-java/klyp-shop/internal/domain/policy"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func InitAddressModule(db *gorm.DB, log *logger.Logger, jwtService *jwt.JWTService) chi.Router {

	userRepository := database.NewUserDB(db)
	addressRepository := database.NewAddressDB(db)
	addressMapper := mapper.NewAddressMapper()
	addressPolicy := policy.NewAddressPolicy()
	addressService := service.NewAddressService(addressRepository, userRepository, log, addressMapper, addressPolicy)
	addressHandler := handler.NewAddressHandler(addressService, log)
	adapter := adapter.NewAdapter(log)
	addressRouter := router.NewAddressRouter(addressHandler, jwtService, userRepository, log, adapter)

	r := chi.NewRouter()
	addressRouter.Handle(r)

	return r
}
