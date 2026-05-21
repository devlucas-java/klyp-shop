package apperrors

import (
	goErrors "errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	pgErrUniqueViolation     = "23505"
	pgErrForeignKeyViolation = "23503"
	pgErrNotNullViolation    = "23502"
	pgErrCheckViolation      = "23514"
	pgErrInvalidText         = "22P02"
)

func HandlePgError(trace string, err error) *DomainError {
	var pgErr *pgconn.PgError

	if goErrors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound(trace+": record not found", err)
	}

	if !goErrors.As(err, &pgErr) {
		return Database(trace+": unexpected database error", err)
	}

	switch pgErr.Code {
	case pgErrUniqueViolation:
		return Conflict(trace+": unique constraint violation", err)
	case pgErrForeignKeyViolation:
		return NotFound(trace+": foreign key constraint violation", err)
	case pgErrNotNullViolation:
		return NotNull(trace+": not null constraint violation", err)
	case pgErrCheckViolation:
		return CheckViolation(trace+": check constraint violation", err)
	case pgErrInvalidText:
		return InvalidText(trace+": invalid text representation", err)
	default:
		return Database(trace+": unexpected database error", err)
	}
}
