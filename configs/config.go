package configs

import (
	"fmt"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type conf struct {
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
}

var cfg *conf

func NewConfig() *conf {
	return cfg
}

// FOR DEVELOPMENT
func InitConfigDev(log *logger.Logger) *conf {

	cfg = &conf{
		log:                log,
		WebServerPort:      "8080",
		DbName:             "klyp_test",
		DbPort:             "5432",
		DbUser:             "postgres",
		DbPassword:         "postgres",
		DbHost:             "localhost",
		JwtSecret:          "test-secret",
		JwtExpireIn:        15,
		JwtRefreshExpireIn: 1440,
	}

	cfg.JwtAccessToken = jwtauth.New("HS256", []byte(cfg.JwtSecret), nil)

	log.Info("config initialized successfully")
	return cfg
}
func InitDBDev(log *logger.Logger) *gorm.DB {

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Error("failed to connect database:", err)
		panic(err)
	}

	err = db.AutoMigrate(
		&entity.Address{},
		&entity.Authority{},
		&entity.BitcoinPayment{},
		&entity.Comment{},
		&entity.Order{},
		&entity.OrderItem{},
		&entity.Product{},
		&entity.Review{},
		&entity.Role{},
		&entity.Seller{},
		&entity.User{},
	)
	if err != nil {
		log.Error("auto migrate failed:", err)
		panic(err)
	}

	log.Info("database initialized successfully")
	return db
}

// FOR TESTING
func InitConfigTest(log *logger.Logger) *conf {

	cfg = &conf{}

	viper.AutomaticEnv()

	err := viper.Unmarshal(cfg)
	if err != nil {

		log.Error("viper unmarshal failed:", err)
		panic(err)
	}

	cfg.JwtAccessToken = jwtauth.New("HS256", []byte(cfg.JwtSecret), nil)

	log.Info("config initialized successfully")
	return cfg
}

func InitDBTest(log *logger.Logger) *gorm.DB {

	cfg := conf{}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("failed to connect database:", err)
		panic(err)
	}

	err = db.AutoMigrate(
		&entity.Address{},
		&entity.Authority{},
		&entity.BitcoinPayment{},
		&entity.Comment{},
		&entity.Order{},
		&entity.OrderItem{},
		&entity.Product{},
		&entity.Review{},
		&entity.Role{},
		&entity.Seller{},
		&entity.User{},
	)
	if err != nil {
		log.Error("auto migrate failed:", err)
		panic(err)
	}

	log.Info("database initialized successfully")
	return db
}

func (c *conf) GetWebServerPort() string {
	return c.WebServerPort
}

func (c *conf) GetJWTSecret() string {
	return c.JwtSecret
}

func (c *conf) GetJWTExpire() int {
	return c.JwtExpireIn
}

func (c *conf) GetTokenAuth() *jwtauth.JWTAuth {
	return c.JwtAccessToken
}
