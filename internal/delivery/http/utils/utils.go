package utils

import (
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

func GetAuth(r *http.Request) (*entity.User, error) {
	auth, ok := r.Context().Value(middleware.AuthKey).(*entity.User)
	if !ok || auth == nil {
		return nil, apperrors.Unauthorized(nil)
	}
	return auth, nil
}
