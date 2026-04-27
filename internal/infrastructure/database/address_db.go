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

type AddressDB struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewAddressDB(db *gorm.DB, log *logger.Logger) repository.AddressRepository {
	return &AddressDB{db: db, log: log}
}

func (a *AddressDB) Create(address *entity.Address) (*entity.Address, error) {
	if err := a.db.Create(address).Error; err != nil {
		a.log.Errorf("db create daddress error: %v", err)
		return nil, fmt.Errorf("create daddress: %w", err)
	}
	return address, nil
}

func (a *AddressDB) Update(address *entity.Address) (*entity.Address, error) {
	if err := a.db.
		Where("id = ?", address.ID).
		Save(address).Error; err != nil {

		a.log.Errorf("db update daddress error (id=%s): %v", address.ID, err)
		return nil, fmt.Errorf("update daddress: %w", err)
	}
	return address, nil
}

func (a *AddressDB) FindByID(id id.UUID) (*entity.Address, error) {
	var address entity.Address

	err := a.db.First(&address, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("Address", err)
		}

		a.log.Errorf("db find daddress by id error (id=%s): %v", id, err)
		return nil, fmt.Errorf("find daddress by id: %w", err)
	}

	return &address, nil
}

func (a *AddressDB) FindByUser(userID id.UUID) ([]*entity.Address, error) {
	var addresses []*entity.Address

	err := a.db.
		Where("user_id = ?", userID).
		Find(&addresses).Error

	if err != nil {
		a.log.Errorf("db find daddress by duser error (user_id=%s): %v", userID, err)
		return nil, fmt.Errorf("find addresses by duser: %w", err)
	}

	return addresses, nil
}

func (a *AddressDB) DeleteByID(id id.UUID) error {
	err := a.db.Delete(&entity.Address{}, "id = ?", id).Error
	if err != nil {
		a.log.Errorf("db delete daddress error (id=%s): %v", id, err)
		return fmt.Errorf("delete daddress: %w", err)
	}

	return nil
}
