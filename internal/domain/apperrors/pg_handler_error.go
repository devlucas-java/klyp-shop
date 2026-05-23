package apperrors

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	pgUniqueViolation     = "23505"
	pgForeignKeyViolation = "23503"
	pgNotNullViolation    = "23502"
	pgCheckViolation      = "23514"
	pgInvalidText         = "22P02"
)

func HandlePgError(resource string, err error) *AppError {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound(resource+" resource was not found", err)
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return Internal(err)
	}

	switch pgErr.Code {
	case pgUniqueViolation:
		return Conflict(resource+" a record with this data already exists", err)
	case pgForeignKeyViolation:
		return Internal(err)
	case pgNotNullViolation:
		return BadRequest(resource+" a required field was not provided", err)
	case pgCheckViolation:
		return BadRequest(resource+" the provided data violates a constraint", err)
	case pgInvalidText:
		return BadRequest(resource+" contains an invalid value", err)
	default:
		return Internal(err)
	}
}
