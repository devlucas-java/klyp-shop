package apperrors

import "net/http"

type AppError struct {
	HTTPStatus  int
	UserMessage string
	Internal    error
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return e.Internal.Error()
	}
	return e.UserMessage
}

func (e *AppError) Unwrap() error {
	return e.Internal
}

func (e *AppError) APIError() (int, string) {
	return e.HTTPStatus, e.UserMessage
}

func NotFound(msg string, cause error) *AppError {
	return &AppError{HTTPStatus: http.StatusNotFound, UserMessage: msg, Internal: cause}
}

func Conflict(msg string, cause error) *AppError {
	return &AppError{HTTPStatus: http.StatusConflict, UserMessage: msg, Internal: cause}
}

func Forbidden(cause error) *AppError {
	return &AppError{HTTPStatus: http.StatusForbidden, UserMessage: "you do not have permission to perform this action", Internal: cause}
}

func Unauthorized(cause error) *AppError {
	return &AppError{HTTPStatus: http.StatusUnauthorized, UserMessage: "authentication required", Internal: cause}
}

func InvalidCredentials(cause error) *AppError {
	return &AppError{HTTPStatus: http.StatusUnauthorized, UserMessage: "invalid credentials", Internal: cause}
}

func BadRequest(msg string, cause error) *AppError {
	return &AppError{HTTPStatus: http.StatusBadRequest, UserMessage: msg, Internal: cause}
}

func Validation(msg string) *AppError {
	return &AppError{HTTPStatus: http.StatusBadRequest, UserMessage: msg, Internal: nil}
}

func InvalidUUID(cause error) *AppError {
	return &AppError{HTTPStatus: http.StatusBadRequest, UserMessage: "the provided identifier is not a valid UUID", Internal: cause}
}

func Unprocessable(msg string, cause error) *AppError {
	return &AppError{HTTPStatus: http.StatusUnprocessableEntity, UserMessage: msg, Internal: cause}
}

func Policy(msg string) *AppError {
	return &AppError{HTTPStatus: http.StatusUnprocessableEntity, UserMessage: msg, Internal: nil}
}

func Internal(cause error) *AppError {
	return &AppError{HTTPStatus: http.StatusInternalServerError, UserMessage: "an unexpected error occurred, please try again later", Internal: cause}
}
