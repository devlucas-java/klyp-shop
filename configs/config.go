package configs

import (
	"fmt"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Conf struct {
	log *logger.Logger

	WebServerPort      string `mapstructure:"WEB_SERVER_PORT"`
	DbName             string `mapstructure:"DB_NAME"`
	DbPort             string `mapstructure:"DB_PORT"`
	DbUser             string `mapstructure:"DB_USER"`
	DbPassword         string `mapstructure:"DB_PASSWORD"`
	DbHost             string `mapstructure:"DB_HOST"`
	DbDriver           string `mapstructure:"DB_DRIVER"`
	JwtSecret          string `mapstructure:"JWT_SECRET"`
	JwtExpireIn        int    `mapstructure:"JWT_EXPIRE_IN"`
	JwtRefreshExpireIn int    `mapstructure:"JWT_REFRESH_EXPIRE_IN"`
	JwtAccessToken     *jwtauth.JWTAuth

	BTCPayBaseURL       string `mapstructure:"BTCPAY_BASE_URL"`
	BTCPayStoreID       string `mapstructure:"BTCPAY_STORE_ID"`
	BTCPayAPIKey        string `mapstructure:"BTCPAY_API_KEY"`
	BTCPayWebhookSecret string `mapstructure:"BTCPAY_WEBHOOK_SECRET"`
}

var cfg *Conf

func NewConfig() *Conf {
	return cfg
}

func InitConfig(log *logger.Logger) *Conf {

	cfg = &Conf{}

	viper.AutomaticEnv()

	viper.BindEnv("WEB_SERVER_PORT")
	viper.BindEnv("DB_DRIVER")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("JWT_EXPIRE_IN")
	viper.BindEnv("JWT_REFRESH_EXPIRE_IN")
	viper.BindEnv("BTCPAY_BASE_URL")
	viper.BindEnv("BTCPAY_STORE_ID")
	viper.BindEnv("BTCPAY_API_KEY")
	viper.BindEnv("BTCPAY_WEBHOOK_SECRET")

	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Error("viper unmarshal failed:", err)
		panic(err)
	}

	log.Infof("DB_HOST: %s, DB_PORT: %s, DB_NAME: %s, DB_USER: %s",
		cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.DbUser)

	cfg.JwtAccessToken = jwtauth.New(
		"HS256",
		[]byte(cfg.JwtSecret),
		nil,
	)

	log.Info("production config initialized successfully")

	return cfg
}

func InitDB(log *logger.Logger) *gorm.DB {

	cfg := NewConfig()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DbHost,
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbName,
		cfg.DbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("failed to connect database:", err)
		panic(err)
	}

	err = db.AutoMigrate(
		&entity.User{},
		&entity.Seller{},
		&entity.Product{},
		&entity.ShoppingCart{},
		&entity.Order{},
		&entity.Address{},
		&entity.BitcoinPayment{},
		&entity.ChatMessage{},
		&entity.Comment{},
		&entity.OrderItem{},
		&entity.Review{},
		&entity.ShoppingCartItem{},
	)

	if err != nil {
		log.Error("auto migrate failed:", err)
		panic(err)
	}

	log.Info("production database initialized successfully")

	return db
}

func (c *Conf) GetWebServerPort() string {
	return c.WebServerPort
}

func (c *Conf) GetDbName() string {
	return c.DbName
}

func (c *Conf) GetDbPort() string {
	return c.DbPort
}

func (c *Conf) GetDbUser() string {
	return c.DbUser
}

func (c *Conf) GetDbPassword() string {
	return c.DbPassword
}

func (c *Conf) GetDbHost() string {
	return c.DbHost
}

func (c *Conf) GetDbDriver() string {
	return c.DbDriver
}

func (c *Conf) GetJWTSecret() string {
	return c.JwtSecret
}

func (c *Conf) GetJWTExpire() int {
	return c.JwtExpireIn
}

func (c *Conf) GetJWTRefreshExpire() int {
	return c.JwtRefreshExpireIn
}

func (c *Conf) GetTokenAuth() *jwtauth.JWTAuth {
	return c.JwtAccessToken
}

func (c *Conf) GetBTCPayBaseURL() string {
	return c.BTCPayBaseURL
}

func (c *Conf) GetBTCPayStoreID() string {
	return c.BTCPayStoreID
}

func (c *Conf) GetBTCPayAPIKey() string {
	return c.BTCPayAPIKey
}

func (c *Conf) GetBTCPayWebhookSecret() string {
	return c.BTCPayWebhookSecret
}
