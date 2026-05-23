package response

import (
	"errors"
	"net/http"

	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi/middleware"
)

type apiErrorInterface interface {
	APIError() (int, string)
}

func ResponseError(w http.ResponseWriter, r *http.Request, err error, log *logger.Logger) {
	requestID := middleware.GetReqID(r.Context())

	var appErr apiErrorInterface
	if errors.As(err, &appErr) {
		status, msg := appErr.APIError()
		logByStatus(log, requestID, status, err)
		ResponseEntity(w, status, map[string]string{"message": msg})
		return
	}

	log.Errorf("[%s] unhandled error (500): %v", requestID, err)
	ResponseEntity(w, http.StatusInternalServerError, map[string]string{
		"message": "an unexpected error occurred, please try again later",
	})
}

func logByStatus(log *logger.Logger, requestID string, status int, err error) {
	switch {
	case status >= 500:
		log.Errorf("[%s] internal error (%d): %v", requestID, status, err)
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		log.Warnf("[%s] auth error (%d): %v", requestID, status, err)
	}
}
