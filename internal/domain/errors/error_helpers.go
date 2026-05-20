package errors

import "errors"

func As(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

func IsAppError(err error) bool {
	return As(err) != nil
}

func IsNotFound(err error) bool {
	var appErr *AppError
	return err != nil && errors.As(err, &appErr) && appErr.Status == 404
}

func IsConflict(err error) bool {
	var appErr *AppError
	return err != nil && errors.As(err, &appErr) && appErr.Status == 409
}

func IsBadRequest(err error) bool {
	var appErr *AppError
	return err != nil && errors.As(err, &appErr) && appErr.Status == 400
}

func IsServerError(err error) bool {
	var appErr *AppError
	return err != nil && errors.As(err, &appErr) && appErr.Status >= 500
}
