package repository

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) (*entity.Order, error)
	Save(ctx context.Context, order *entity.Order) (*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) (*entity.Order, error)
	Updates(ctx context.Context, order *entity.Order) (*entity.Order, error)
	FindByID(ctx context.Context, id id.UUID) (*entity.Order, error)
	FindByUserIDPaginated(ctx context.Context, userID id.UUID, page, size int, status string) ([]*entity.Order, int64, error)
	FindAllPaginated(ctx context.Context, page, size int, status string) ([]*entity.Order, int64, error)
	FindBySellerIDPaginated(ctx context.Context, sellerID id.UUID, page, size int, status string) ([]*entity.Order, int64, error)
	DeleteByID(ctx context.Context, id id.UUID) error
}
