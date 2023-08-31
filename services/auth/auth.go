package auth

import (
	"auth/conf"
	"auth/log"
	
	"gorm.io/gorm"

	//"time"

	//"gorm.io/gorm"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	//"github.com/google/uuid"
)

type Service struct {
	Log   *log.Logger
	AuthStore *AuthStore
	Validator *validator.Validate
}

type AuthStore struct {
	log  *log.Logger
	conf *conf.Config
	db   *gorm.DB
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MyCustomClaims struct {
	Username string `json:"username"`
	Scopes string `json:"scopes"`
	jwt.RegisteredClaims
}

func NewService(log *log.Logger, conf *conf.Config, db *gorm.DB, validator *validator.Validate) *Service {
	return &Service{
		Log:   log,
		AuthStore: NewStore(log, conf, db),
		Validator: validator,
	}
}

func NewStore(log *log.Logger, conf *conf.Config, db *gorm.DB) *AuthStore {
	return &AuthStore{
		log:  log,
		// conf: conf,
		db:   db,
	}
}
