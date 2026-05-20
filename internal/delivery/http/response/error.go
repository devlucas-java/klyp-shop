package response

import (
	"errors"
	"net/http"

	domainErrors "github.com/devlucas-java/klyp-shop/internal/domain/errors"
)

func ResponseError(w http.ResponseWriter, err error) {
	var appErr *domainErrors.AppError
	if errors.As(err, &appErr) {
		ResponseEntity(w, appErr.StatusCode(), appErr)
		return
	}

	ResponseEntity(w, http.StatusInternalServerError, domainErrors.New("an unexpected error occurred", http.StatusInternalServerError, err))
}
