package errors

func ErrInvalidCredentials(err error) *AppError {
	return New(
		"Invalid email/username or password",
		401,
		err,
	)
}

func ErrInvalidRole(message string, err error) *AppError {
	return New(
		message,
		400,
		err,
	)
}

func ErrInvalidUUID(err error) *AppError {
	return New(
		"Invalid UUID provided",
		400,
		err,
	)
}

func ErrNotFound(resource string, err error) *AppError {
	return New(
		resource+" not found",
		404,
		err,
	)
}
func ErrConflict(resource string, err error) *AppError {
	return New(
		resource+" already exists",
		409,
		err,
	)
}
func ErrUnauthorized(err error) *AppError {
	return New(
		"Unauthorized",
		401,
		err,
	)
}
func ErrForbidden(err error) *AppError {
	return New(
		"Forbidden",
		403,
		err,
	)
}
func ErrBadRequest(message string, err error) *AppError {
	return New(
		message,
		400,
		err,
	)
}
func ErrInvalidPayload(err error) *AppError {
	return New(
		"invalid request payload",
		400,
		err,
	)
}
func ErrUnprocessable(message string, err error) *AppError {
	return New(
		message,
		422,
		err,
	)
}
func ErrDatabase(message string, err error) *AppError {
	return New(
		message,
		500,
		err,
	)
}
func ErrInternal(message string, err error) *AppError {
	return New(
		message,
		500,
		err,
	)
}
