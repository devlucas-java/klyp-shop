package errors

import "fmt"

type AppError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Status, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Status, e.Message)
}

func (e *AppError) StatusCode() int {
	return e.Status
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(message string, status int, err error) *AppError {
	return &AppError{
		Message: message,
		Status:  status,
		Err:     err,
	}
}

func Wrap(message string, status int, err error) *AppError {
	return New(message, status, err)
}
