package middleware

import (
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

func RoleMiddleware(roles []enums.Role, log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth, ok := r.Context().Value(AuthKey).(*entity.User)
			if !ok || auth == nil {
				response.ResponseError(w, r, apperrors.Unauthorized(nil), log)
				return
			}

			for _, required := range roles {
				for _, userRole := range auth.Roles {
					if required == userRole {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			response.ResponseError(w, r, apperrors.Forbidden(nil), log)
		})
	}
}
