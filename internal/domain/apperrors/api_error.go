package apperrors

import (
	"errors"
	"net/http"
)

type APIError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(message string, status int) *APIError {
	return &APIError{Message: message, Status: status}
}

func ToAPIError(err error) *APIError {
	var domainErr *DomainError
	if !errors.As(err, &domainErr) {
		return NewAPIError("an unexpected error occurred", http.StatusInternalServerError)
	}

	switch domainErr.Kind {
	case KindNotFound:
		return NewAPIError(domainErr.Message, http.StatusNotFound)
	case KindConflict:
		return NewAPIError(domainErr.Message, http.StatusConflict)
	case KindForbidden:
		return NewAPIError("you do not have permission to perform this action", http.StatusForbidden)
	case KindUnauthorized:
		return NewAPIError("you do not have permission to perform this action", http.StatusForbidden)
	case KindPolicy:
		return NewAPIError(domainErr.Message, http.StatusUnprocessableEntity)
	case KindUnprocessable:
		return NewAPIError(domainErr.Message, http.StatusUnprocessableEntity)
	case KindBadRequest, KindValidation, KindNotNull, KindCheckViolation, KindInvalidText, KindInvalidUUID:
		return NewAPIError(domainErr.Message, http.StatusBadRequest)
	case KindDatabase, KindInternal:
		return NewAPIError("an unexpected error occurred", http.StatusInternalServerError)
	default:
		return NewAPIError("an unexpected error occurred", http.StatusInternalServerError)
	}
}
