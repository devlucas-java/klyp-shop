package utils

import (
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

func GetAuth(r *http.Request) (*entity.User, error) {
	auth, ok := r.Context().Value(middleware.AuthKey).(*entity.User)
	if !ok {
		return nil, apperrors.Unauthorized("unauthorized", nil)
	}
	return auth, nil
}
