package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type BitcoinPaymentRepository interface {
	Create(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error)
	Save(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error)
	Update(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error)
	Updates(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error)
	FindByID(id id.UUID) (*entity.BitcoinPayment, error)
	FindByOrderID(orderID id.UUID) (*entity.BitcoinPayment, error)
	FindByTxHash(txHash string) (*entity.BitcoinPayment, error)
	DeleteByID(id id.UUID) error
}
