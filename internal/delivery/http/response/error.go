package response

import (
	"errors"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

func ResponseError(w http.ResponseWriter, err error, log *logger.Logger) {
	var apiErr *apperrors.APIError
	if errors.As(err, &apiErr) {
		ResponseEntity(w, apiErr.Status, apiErr)
		return
	}

	mapped := apperrors.ToAPIError(err)

	logByKind(log, err, mapped)

	ResponseEntity(w, mapped.Status, mapped)
}

func logByKind(log *logger.Logger, original error, mapped *apperrors.APIError) {
	var domainErr *apperrors.DomainError
	if !errors.As(original, &domainErr) {
		log.Errorf("unhandled error: %v", original)
		return
	}

	switch domainErr.Kind {
	case apperrors.KindForbidden, apperrors.KindUnauthorized:
		log.Warnf("[%s] %v", domainErr.Kind, original)

	case apperrors.KindDatabase, apperrors.KindInternal:
		log.Errorf("[%s] %v", domainErr.Kind, original)

	case apperrors.KindNotFound, apperrors.KindConflict,
		apperrors.KindUnprocessable, apperrors.KindBadRequest,
		apperrors.KindInvalidUUID, apperrors.KindNotNull,
		apperrors.KindCheckViolation, apperrors.KindInvalidText:
	}
}
