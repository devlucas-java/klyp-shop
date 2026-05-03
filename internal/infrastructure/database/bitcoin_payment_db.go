package database

import (
	"errors"
	"fmt"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"gorm.io/gorm"
)

type BitcoinPaymentDB struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewBitcoinPaymentDB(db *gorm.DB, log *logger.Logger) repository.BitcoinPaymentRepository {
	return &BitcoinPaymentDB{db: db, log: log}
}

func (b *BitcoinPaymentDB) Create(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error) {
	if err := b.db.Create(payment).Error; err != nil {
		b.log.Errorf("BitcoinPaymentDB.Create: %v", err)
		return nil, fmt.Errorf("failed to create bitcoin payment: %w", err)
	}
	return payment, nil
}

func (b *BitcoinPaymentDB) Save(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error) {
	if err := b.db.Where("id = ?", payment.ID).Save(payment).Error; err != nil {
		b.log.Errorf("BitcoinPaymentDB.Save %s: %v", payment.ID, err)
		return nil, fmt.Errorf("failed to save bitcoin payment: %w", err)
	}
	return payment, nil
}

func (b *BitcoinPaymentDB) Updates(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error) {
	if err := b.db.Model(payment).Where("id = ?", payment.ID).Updates(payment).Error; err != nil {
		b.log.Errorf("BitcoinPaymentDB.Updates %s: %v", payment.ID, err)
		return nil, fmt.Errorf("failed to update bitcoin payment: %w", err)
	}
	return payment, nil
}

func (b *BitcoinPaymentDB) Update(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error) {
	return b.Save(payment)
}

func (b *BitcoinPaymentDB) FindByID(paymentID id.UUID) (*entity.BitcoinPayment, error) {
	var payment entity.BitcoinPayment
	err := b.db.First(&payment, "id = ?", paymentID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("BitcoinPayment", err)
		}
		b.log.Errorf("BitcoinPaymentDB.FindByID %s: %v", paymentID, err)
		return nil, fmt.Errorf("failed to find bitcoin payment: %w", err)
	}
	return &payment, nil
}

func (b *BitcoinPaymentDB) FindByOrderID(orderID id.UUID) (*entity.BitcoinPayment, error) {
	var payment entity.BitcoinPayment
	err := b.db.Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("BitcoinPayment", err)
		}
		b.log.Errorf("BitcoinPaymentDB.FindByOrderID %s: %v", orderID, err)
		return nil, fmt.Errorf("failed to find bitcoin payment: %w", err)
	}
	return &payment, nil
}

func (b *BitcoinPaymentDB) FindByTxHash(txHash string) (*entity.BitcoinPayment, error) {
	var payment entity.BitcoinPayment
	err := b.db.Where("tx_hash = ?", txHash).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("BitcoinPayment", err)
		}
		b.log.Errorf("BitcoinPaymentDB.FindByTxHash %s: %v", txHash, err)
		return nil, fmt.Errorf("failed to find bitcoin payment: %w", err)
	}
	return &payment, nil
}

func (b *BitcoinPaymentDB) DeleteByID(paymentID id.UUID) error {
	if err := b.db.Delete(&entity.BitcoinPayment{}, "id = ?", paymentID).Error; err != nil {
		b.log.Errorf("BitcoinPaymentDB.DeleteByID %s: %v", paymentID, err)
		return fmt.Errorf("failed to delete bitcoin payment: %w", err)
	}
	return nil
}
