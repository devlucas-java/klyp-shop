package configs

import (
	"fmt"

	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type conf struct {
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
func InitConfigDev() *conf {

	cfg = &conf{
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

	return cfg
}
func InitDBDev() *gorm.DB {

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate()
	if err != nil {
		panic(err)
	}

	return db
}

// FOR TESTING

func InitConfigTest() *conf {

	cfg = &conf{}

	viper.AutomaticEnv()

	err := viper.Unmarshal(cfg)
	if err != nil {
		panic(err)
	}

	cfg.JwtAccessToken = jwtauth.New("HS256", []byte(cfg.JwtSecret), nil)

	return cfg
}

func InitDBTest() *gorm.DB {

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
		panic(err)
	}

	err = db.AutoMigrate()
	if err != nil {
		panic(err)
	}

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
