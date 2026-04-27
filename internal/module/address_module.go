package module

import (
	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/router"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func InitAddressModule(db *gorm.DB, log *logger.Logger, jwtService *jwt.JWTService) chi.Router {

	userRepository := database.NewUserDB(db, log)
	addressRepository := database.NewAddressDB(db, log)
	addressMapper := mapper.NewAddressMapper()
	addressService := service.NewAddressService(addressRepository, log, addressMapper, userRepository)
	addressHandler := handler.NewAddressHandler(addressService, log)
	addressRouter := router.NewAddressRouter(addressHandler, log, jwtService, userRepository)

	r := chi.NewRouter()
	addressRouter.Handle(r)

	return r
}
