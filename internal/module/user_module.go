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

func InitUserModule(db *gorm.DB, log *logger.Logger, jwtService *jwt.JWTService) chi.Router {

	userRepository := database.NewUserDB(db, log)
	userMapper := mapper.NewUserMapper()
	userService := service.NewUserService(userRepository, log, userMapper)
	userHandler := handler.NewUserHandler(userService, log)
	UserRouter := router.NewUserRouter(jwtService, userHandler, log, userRepository)

	r := chi.NewRouter()
	UserRouter.RegisterUserRoutes(r)

	return r
}
