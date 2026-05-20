package errors

import (
	goErrors "errors"
	"strings"

	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgErrUniqueViolation     = "23505"
	pgErrForeignKeyViolation = "23503"
	pgErrNotNullViolation    = "23502"
	pgErrCheckViolation      = "23514"
	pgErrInvalidText         = "22P02"
)

func HandlePgError(log *logger.Logger, err error, genericMsg string) error {
	var pgErr *pgconn.PgError

	if !goErrors.As(err, &pgErr) {
		log.Errorf(
			"pg error detail=%v message=%s",
			err.Error(),
			genericMsg,
		)
		return ErrDatabase(genericMsg, err)
	}

	field := extractField(pgErr.Detail)

	switch pgErr.Code {

	case pgErrUniqueViolation:
		return ErrConflict(
			field+" already used",
			err,
		)

	case pgErrForeignKeyViolation:
		return ErrBadRequest(
			field+" does not exist",
			err,
		)

	case pgErrNotNullViolation:
		return ErrBadRequest(
			field+" is required",
			err,
		)

	case pgErrCheckViolation:
		return ErrBadRequest(
			field+" is invalid",
			err,
		)

	case pgErrInvalidText:
		return ErrBadRequest(
			field+" has invalid format",
			err,
		)

	default:

		log.Errorf(
			"pg error code=%s field=%s detail=%s message=%s",
			pgErr.Code,
			field,
			pgErr.Detail,
			pgErr.Message,
		)
		return ErrDatabase(genericMsg, err)
	}
}

func extractField(detail string) string {
	start := strings.Index(detail, "(")
	end := strings.Index(detail, ")")

	if start == -1 || end == -1 || start >= end {
		return "resource"
	}

	return strings.TrimSpace(detail[start+1 : end])
}
