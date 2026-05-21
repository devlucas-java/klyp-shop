package apperrors

import "fmt"

type ErrorKind string

const (
	KindNotFound       ErrorKind = "not_found"
	KindConflict       ErrorKind = "conflict"
	KindForbidden      ErrorKind = "forbidden"
	KindUnprocessable  ErrorKind = "unprocessable"
	KindBadRequest     ErrorKind = "bad_request"
	KindUnauthorized   ErrorKind = "unauthorized"
	KindDatabase       ErrorKind = "database"
	KindInternal       ErrorKind = "internal"
	KindNotNull        ErrorKind = "not_null"
	KindCheckViolation ErrorKind = "check_violation"
	KindInvalidText    ErrorKind = "invalid_text"
	KindInvalidUUID    ErrorKind = "invalid_uuid"
	KindValidation     ErrorKind = "validation"
	KindPolicy         ErrorKind = "policy"
)

type DomainError struct {
	Kind    ErrorKind
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}

func newDomainError(kind ErrorKind, message string, cause error) *DomainError {
	return &DomainError{Kind: kind, Message: message, Cause: cause}
}

func NotFound(message string, cause error) *DomainError {
	return newDomainError(KindNotFound, message, cause)
}

func Conflict(message string, cause error) *DomainError {
	return newDomainError(KindConflict, message, cause)
}

func Forbidden(message string, cause error) *DomainError {
	return newDomainError(KindForbidden, message, cause)
}

func Unprocessable(message string, cause error) *DomainError {
	return newDomainError(KindUnprocessable, message, cause)
}

func BadRequest(message string, cause error) *DomainError {
	return newDomainError(KindBadRequest, message, cause)
}

func Unauthorized(message string, cause error) *DomainError {
	return newDomainError(KindUnauthorized, message, cause)
}

func NotNull(message string, cause error) *DomainError {
	return newDomainError(KindNotNull, message, cause)
}

func CheckViolation(message string, cause error) *DomainError {
	return newDomainError(KindCheckViolation, message, cause)
}

func InvalidText(message string, cause error) *DomainError {
	return newDomainError(KindInvalidText, message, cause)
}

func InvalidUUID(message string, cause error) *DomainError {
	return newDomainError(KindInvalidUUID, message, cause)
}

func Database(message string, cause error) *DomainError {
	return newDomainError(KindDatabase, message, cause)
}

func Internal(message string, cause error) *DomainError {
	return newDomainError(KindInternal, message, cause)
}

func Validation(message string) *DomainError {
	return newDomainError(KindValidation, message, nil)
}

func Policy(message string) *DomainError {
	return newDomainError(KindPolicy, message, nil)
}
