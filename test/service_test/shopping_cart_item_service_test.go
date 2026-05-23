package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	cartDTO "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/cart"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/devlucas-java/klyp-shop/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newShoppingCartItemService(cartRepo *mocks.ShoppingCartRepositoryMock, itemRepo *mocks.ShoppingCartItemRepositoryMock, productRepo *mocks.ProductRepositoryMock) *service.ShoppingCartItemService {
	return service.NewShoppingCartItemService(
		logger.NewLogger(logger.TRACE),
		cartRepo,
		itemRepo,
		productRepo,
		mapper.NewShoppingCartMapper(),
	)
}

// helpers

func newCartWithItems(userID id.UUID, items ...*entity.ShoppingCartItem) *entity.ShoppingCart {
	c := &entity.ShoppingCart{
		ID:     id.NewUUID(),
		UserID: userID,
		Items:  items,
	}
	c.RecalculateTotal()
	return c
}

// ── AddItem ──────────────────────────────────────────────────────────────────

func TestShoppingCartItemService_AddItem_NewProduct(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	itemRepo := new(mocks.ShoppingCartItemRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newShoppingCartItemService(cartRepo, itemRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	product := &entity.Product{ID: id.NewUUID(), PriceBTC: 1000}
	cart := newCartWithItems(user.ID)

	productRepo.On("FindByID", product.ID).Return(product, nil)
	cartRepo.On("FindByUserID", user.ID).Return(cart, nil)
	cartRepo.On("Save", mock.AnythingOfType("*entity.ShoppingCart")).Return(cart, nil)

	req := &cartDTO.AddShoppingCartItemRequest{ProductID: product.ID.String(), Quantity: 2}
	res, err := svc.AddItem(user, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, user.ID.String(), res.UserID)
	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestShoppingCartItemService_AddItem_DuplicateProduct_IncrementsQuantity(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	itemRepo := new(mocks.ShoppingCartItemRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newShoppingCartItemService(cartRepo, itemRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	product := &entity.Product{ID: id.NewUUID(), PriceBTC: 500}
	cartID := id.NewUUID()
	existingItem := &entity.ShoppingCartItem{
		ID:        id.NewUUID(),
		CartID:    cartID,
		ProductID: product.ID,
		Quantity:  1,
		PriceBTC:  500,
	}
	cart := &entity.ShoppingCart{ID: cartID, UserID: user.ID, Items: []*entity.ShoppingCartItem{existingItem}}
	cart.RecalculateTotal()

	productRepo.On("FindByID", product.ID).Return(product, nil)
	cartRepo.On("FindByUserID", user.ID).Return(cart, nil)
	cartRepo.On("Save", mock.AnythingOfType("*entity.ShoppingCart")).Return(cart, nil)

	req := &cartDTO.AddShoppingCartItemRequest{ProductID: product.ID.String(), Quantity: 3}
	res, err := svc.AddItem(user, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	// quantity should be 1 + 3 = 4
	assert.Equal(t, 4, existingItem.Quantity)
	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestShoppingCartItemService_AddItem_InvalidQuantity(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	itemRepo := new(mocks.ShoppingCartItemRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newShoppingCartItemService(cartRepo, itemRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	req := &cartDTO.AddShoppingCartItemRequest{ProductID: id.NewUUID().String(), Quantity: 0}

	_, err := svc.AddItem(user, req)

	assert.Error(t, err)
	cartRepo.AssertNotCalled(t, "FindByUserID")
	productRepo.AssertNotCalled(t, "FindByID")
}

func TestShoppingCartItemService_AddItem_ProductNotFound(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	itemRepo := new(mocks.ShoppingCartItemRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newShoppingCartItemService(cartRepo, itemRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	ghostID := id.NewUUID()

	productRepo.On("FindByID", ghostID).Return(nil, apperrors.NotFound("product", nil))

	req := &cartDTO.AddShoppingCartItemRequest{ProductID: ghostID.String(), Quantity: 1}
	_, err := svc.AddItem(user, req)

	assert.Error(t, err)
	cartRepo.AssertNotCalled(t, "FindByUserID")
	productRepo.AssertExpectations(t)
}

// ── UpdateItem ───────────────────────────────────────────────────────────────

func TestShoppingCartItemService_UpdateItem(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	itemRepo := new(mocks.ShoppingCartItemRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newShoppingCartItemService(cartRepo, itemRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	cartID := id.NewUUID()
	item := &entity.ShoppingCartItem{
		ID:       id.NewUUID(),
		CartID:   cartID,
		Quantity: 2,
		PriceBTC: 300,
	}
	cart := &entity.ShoppingCart{ID: cartID, UserID: user.ID, Items: []*entity.ShoppingCartItem{item}}
	cart.RecalculateTotal()

	cartRepo.On("FindByUserID", user.ID).Return(cart, nil)
	cartRepo.On("Save", mock.AnythingOfType("*entity.ShoppingCart")).Return(cart, nil)

	req := &cartDTO.UpdateShoppingCartItemRequest{Quantity: 5}
	res, err := svc.UpdateItem(user, item.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 5, item.Quantity)
	assert.Equal(t, int64(1500), cart.TotalBTC) // 300 * 5
	cartRepo.AssertExpectations(t)
}

func TestShoppingCartItemService_UpdateItem_InvalidQuantity(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	itemRepo := new(mocks.ShoppingCartItemRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newShoppingCartItemService(cartRepo, itemRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	req := &cartDTO.UpdateShoppingCartItemRequest{Quantity: -1}

	_, err := svc.UpdateItem(user, id.NewUUID(), req)

	assert.Error(t, err)
	cartRepo.AssertNotCalled(t, "FindByUserID")
}

func TestShoppingCartItemService_UpdateItem_ItemNotInCart(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	itemRepo := new(mocks.ShoppingCartItemRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newShoppingCartItemService(cartRepo, itemRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	cart := newCartWithItems(user.ID) // empty cart

	cartRepo.On("FindByUserID", user.ID).Return(cart, nil)

	req := &cartDTO.UpdateShoppingCartItemRequest{Quantity: 3}
	_, err := svc.UpdateItem(user, id.NewUUID(), req)

	assert.Error(t, err)
	cartRepo.AssertExpectations(t)
}

// ── RemoveItem ───────────────────────────────────────────────────────────────

func TestShoppingCartItemService_RemoveItem(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	itemRepo := new(mocks.ShoppingCartItemRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newShoppingCartItemService(cartRepo, itemRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	cartID := id.NewUUID()
	item := &entity.ShoppingCartItem{
		ID:       id.NewUUID(),
		CartID:   cartID,
		Quantity: 1,
		PriceBTC: 200,
	}
	cart := &entity.ShoppingCart{ID: cartID, UserID: user.ID, Items: []*entity.ShoppingCartItem{item}}
	cart.RecalculateTotal()
	savedCart := &entity.ShoppingCart{ID: cartID, UserID: user.ID, Items: []*entity.ShoppingCartItem{}}

	cartRepo.On("FindByUserID", user.ID).Return(cart, nil)
	itemRepo.On("DeleteByID", item.ID).Return(nil)
	cartRepo.On("Save", mock.AnythingOfType("*entity.ShoppingCart")).Return(savedCart, nil)

	res, err := svc.RemoveItem(user, item.ID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(0), cart.TotalBTC) // recalculated after removal
	cartRepo.AssertExpectations(t)
	itemRepo.AssertExpectations(t)
}

func TestShoppingCartItemService_RemoveItem_ItemNotFound(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	itemRepo := new(mocks.ShoppingCartItemRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newShoppingCartItemService(cartRepo, itemRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	cart := newCartWithItems(user.ID) // empty cart

	cartRepo.On("FindByUserID", user.ID).Return(cart, nil)

	_, err := svc.RemoveItem(user, id.NewUUID())

	assert.Error(t, err)
	cartRepo.AssertNotCalled(t, "Save")
	itemRepo.AssertNotCalled(t, "DeleteByID")
	cartRepo.AssertExpectations(t)
	itemRepo.AssertExpectations(t)
}
