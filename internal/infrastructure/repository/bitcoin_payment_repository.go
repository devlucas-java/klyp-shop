package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type BitcoinPaymentRepository interface {
	Save(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error)
	Updates(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error)
	FindByOrderID(orderID id.UUID) (*entity.BitcoinPayment, error)
	FindByTxHash(txHash string) (*entity.BitcoinPayment, error)
}
