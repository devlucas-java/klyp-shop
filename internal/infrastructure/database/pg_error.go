package database

import (
	"errors"

	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgErrUniqueViolation     = "23505"
	pgErrForeignKeyViolation = "23503"
	pgErrNotNullViolation    = "23502"
	pgErrCheckViolation      = "23514"
	pgErrInvalidText         = "22P02"
)

func handlePgError(err error, genericMsg string) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {

		switch pgErr.Code {

		case pgErrUniqueViolation:
			return domainErr.ErrConflict("resource already exists", err)

		case pgErrForeignKeyViolation:
			return domainErr.ErrNotFound("related resource not found", err)

		case pgErrNotNullViolation:
			return domainErr.ErrBadRequest("missing required field", err)

		case pgErrCheckViolation:
			return domainErr.ErrBadRequest("validation failed", err)

		case pgErrInvalidText:
			return domainErr.ErrBadRequest("invalid input format", err)
		}
	}

	return domainErr.ErrDatabase(genericMsg, err)
}
