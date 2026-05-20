package middleware

import (
	"context"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/jwtauth"
)

type contextKey string

const (
	AuthKey contextKey = "auth_context"
	JTIKey  contextKey = "jti_context"
)

func JwtMiddleware(
	jwtService *jwt.JWTService,
	log *logger.Logger,
	userRepository repository.UserRepository,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := jwtauth.TokenFromHeader(r)

			claims, err := jwtService.Validate(token)
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userID, err := id.Parse(userIDStr)
			if err != nil {
				http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
				return
			}

			email, _ := claims["email"].(string)
			jti, _ := claims["jti"].(string)

			user, err := userRepository.FindByID(userID)
			if err != nil || user == nil {
				log.Errorf("JwtMiddleware: user not found %s: %v", userID, err)

				http.Error(w, "Unauthorize", http.StatusUnauthorized)
				return
			}

			var roles []enums.Role
			if rolesClaim, ok := claims["roles"].([]interface{}); ok {
				for _, claim := range rolesClaim {
					if roleStr, ok := claim.(string); ok {
						roles = append(roles, enums.Role(roleStr))
					}
				}
			}

			auth := &entity.User{
				ID:    userID,
				Email: email,
				Roles: roles,
			}

			ctx := context.WithValue(r.Context(), AuthKey, auth)
			ctx = context.WithValue(ctx, JTIKey, jti)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
