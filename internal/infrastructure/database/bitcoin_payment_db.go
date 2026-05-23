package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

type BitcoinPaymentDB struct {
	db *gorm.DB
}

func NewBitcoinPaymentDB(db *gorm.DB) repository.BitcoinPaymentRepository {
	return &BitcoinPaymentDB{db: db}
}

func (b *BitcoinPaymentDB) Create(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error) {
	if err := b.db.WithContext(context.Background()).Create(payment).Error; err != nil {
		return nil, apperrors.HandlePgError("payment", err)
	}
	return payment, nil
}

func (b *BitcoinPaymentDB) Save(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error) {
	if err := b.db.WithContext(context.Background()).Where("id = ?", payment.ID).Save(payment).Error; err != nil {
		return nil, apperrors.HandlePgError("payment", err)
	}
	return payment, nil
}

func (b *BitcoinPaymentDB) Updates(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error) {
	if err := b.db.WithContext(context.Background()).Model(payment).Where("id = ?", payment.ID).Updates(payment).Error; err != nil {
		return nil, apperrors.HandlePgError("payment", err)
	}
	return payment, nil
}

func (b *BitcoinPaymentDB) Update(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error) {
	return b.Save(payment)
}

func (b *BitcoinPaymentDB) FindByID(paymentID id.UUID) (*entity.BitcoinPayment, error) {
	var payment entity.BitcoinPayment
	err := b.db.WithContext(context.Background()).First(&payment, "id = ?", paymentID).Error
	if err != nil {
		return nil, apperrors.HandlePgError("payment", err)
	}
	return &payment, nil
}

func (b *BitcoinPaymentDB) FindByOrderID(orderID id.UUID) (*entity.BitcoinPayment, error) {
	var payment entity.BitcoinPayment
	err := b.db.WithContext(context.Background()).Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, apperrors.HandlePgError("payment", err)
	}
	return &payment, nil
}

func (b *BitcoinPaymentDB) FindByTxHash(txHash string) (*entity.BitcoinPayment, error) {
	var payment entity.BitcoinPayment
	err := b.db.WithContext(context.Background()).Where("tx_hash = ?", txHash).First(&payment).Error
	if err != nil {
		return nil, apperrors.HandlePgError("payment", err)
	}
	return &payment, nil
}

func (b *BitcoinPaymentDB) DeleteByID(paymentID id.UUID) error {
	if err := b.db.WithContext(context.Background()).Delete(&entity.BitcoinPayment{}, "id = ?", paymentID).Error; err != nil {
		return apperrors.HandlePgError("payment", err)
	}
	return nil
}
