package database

import (
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

func (b *BitcoinPaymentDB) Save(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error) {
	if err := b.db.Create(payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

func (b *BitcoinPaymentDB) Updates(payment *entity.BitcoinPayment) (*entity.BitcoinPayment, error) {
	err := b.db.
		Model(&entity.BitcoinPayment{}).
		Where("id = ?", payment.ID).
		Updates(payment).Error
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (b *BitcoinPaymentDB) FindByOrderID(orderID id.UUID) (*entity.BitcoinPayment, error) {
	var payment entity.BitcoinPayment
	err := b.db.
		Where("order_id = ?", orderID).
		First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (b *BitcoinPaymentDB) FindByTxHash(txHash string) (*entity.BitcoinPayment, error) {
	var payment entity.BitcoinPayment
	err := b.db.
		Where("tx_hash = ?", txHash).
		First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
