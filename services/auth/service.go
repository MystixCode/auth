package auth

import (
	"auth/conf"
	"auth/log"
	
	"gorm.io/gorm"
	"github.com/go-playground/validator/v10"
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
