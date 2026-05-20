package database

import (
	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
)

func handlePgError(err error, genericMsg string) error {
	return domainErr.HandlePgError(err, genericMsg)
}
