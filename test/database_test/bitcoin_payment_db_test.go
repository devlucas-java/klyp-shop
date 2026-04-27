package database_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbBitcoin *gorm.DB
var bitcoinRepo *database.BitcoinPaymentDB

func setupBitcoinDB(t *testing.T) {
	var err error

	dbBitcoin, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = dbBitcoin.AutoMigrate(
		&entity.User{},
		&entity.Address{},
		&entity.Seller{},
		&entity.Product{},
		&entity.Order{},
		&entity.OrderItem{},
		&entity.BitcoinPayment{},
	)
	if err != nil {
		t.Fatal(err)
	}

	bitcoinRepo = database.NewBitcoinPaymentDB(dbBitcoin).(*database.BitcoinPaymentDB)
}

func createOrderForPayment(t *testing.T) *entity.Order {
	user := &entity.User{
		ID:       id.NewUUID(),
		Name:     "btc-user",
		Email:    "btcuser@test.com",
		Username: "btcuser",
		Password: "hash",
	}
	dbBitcoin.Create(user)

	addr := &entity.Address{
		ID:       id.NewUUID(),
		UserID:   user.ID,
		Street:   "BTC Street",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "00000",
	}
	dbBitcoin.Create(addr)

	order := entity.NewOrder(user.ID, addr.ID, nil)
	dbBitcoin.Create(order)

	return order
}

func TestSaveBitcoinPayment(t *testing.T) {
	setupBitcoinDB(t)

	order := createOrderForPayment(t)

	payment := entity.NewBitcoinPayment(order.ID, "bc1qxyz123", 0.005)

	res, err := bitcoinRepo.Save(payment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, payment.WalletAddress, res.WalletAddress)
	assert.Equal(t, payment.AmountBTC, res.AmountBTC)
	assert.Equal(t, entity.PaymentStatusPending, res.Status)
	assert.Equal(t, order.ID, res.OrderID)
}

func TestFindBitcoinPaymentByOrderID(t *testing.T) {
	setupBitcoinDB(t)

	order := createOrderForPayment(t)

	payment := entity.NewBitcoinPayment(order.ID, "bc1qabc456", 0.01)
	dbBitcoin.Create(payment)

	found, err := bitcoinRepo.FindByOrderID(order.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, payment.ID, found.ID)
	assert.Equal(t, payment.WalletAddress, found.WalletAddress)
}

func TestFindBitcoinPaymentByOrderID_NotFound(t *testing.T) {
	setupBitcoinDB(t)

	_, err := bitcoinRepo.FindByOrderID(id.NewUUID())
	assert.Error(t, err)
}

func TestFindBitcoinPaymentByTxHash(t *testing.T) {
	setupBitcoinDB(t)

	order := createOrderForPayment(t)

	payment := entity.NewBitcoinPayment(order.ID, "bc1qdef789", 0.02)
	payment.TxHash = "txhash_abc123"
	dbBitcoin.Create(payment)

	found, err := bitcoinRepo.FindByTxHash("txhash_abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, payment.ID, found.ID)
	assert.Equal(t, "txhash_abc123", found.TxHash)
}

func TestFindBitcoinPaymentByTxHash_NotFound(t *testing.T) {
	setupBitcoinDB(t)

	_, err := bitcoinRepo.FindByTxHash("nonexistent_hash")
	assert.Error(t, err)
}

func TestUpdateBitcoinPayment(t *testing.T) {
	setupBitcoinDB(t)

	order := createOrderForPayment(t)

	payment := entity.NewBitcoinPayment(order.ID, "bc1qghi000", 0.03)
	dbBitcoin.Create(payment)

	payment.Confirm("txhash_confirmed")

	res, err := bitcoinRepo.Updates(payment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, entity.PaymentStatusConfirmed, res.Status)
	assert.Equal(t, "txhash_confirmed", res.TxHash)
}
