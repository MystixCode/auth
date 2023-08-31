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



type App struct {
	ID            	int    `json:"id" gorm:"primaryKey"`
	AppName       	string `json:"app_name"`
	AppURI        	string `json:"app_uri"`
	RedirectURI   	string `json:"redirect_uri"`
	ClientType    	string `json:"client_type"` //public or confidential. maybe create a table for it
	Alg				string `json:"alg"`
	ClientID		string `json:"client_id"`
	CreatedAt     	int64  `json:"created_at"`
	UpdatedAt     	int64  `json:"updated_at"`

	Keys []Key `json:"keys" gorm:"foreignKey:AppID;onDelete:CASCADE"` // Many keys belong to one app
}

type AppInput struct {
	AppName  		string `json:"app_name" validate:"required"`
	AppURI      	string `json:"app_uri"`
	Alg				string `json:"alg" validate:"required"`
	RedirectURI   	string `json:"redirect_uri"`
	ClientType    	string `json:"client_type"`
}


type Key struct {
	ID         	int		`json:"id" gorm:"primaryKey"`
	AppID		int		`json:"app_id" gorm:"not null"`
	Alg			string	`json:"alg" gorm:"not null"`
	CreatedAt	int64	`json:"created_at" gorm:"not null"`
}

type KeyInput struct {
	AppID 	int		`json:"app_id" validate:"required,number"`
	Alg		string	`json:"alg" validate:"required,alphanum"`
}


type User struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	UserName  string `json:"user_name" gorm:"not null"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Hash      string `json:"hash" gorm:"not null"`
	CreatedAt int64  `json:"created_at" gorm:"not null"`
	UpdatedAt int64  `json:"updated_at"`
}

type UserInput struct {
	UserName  string `json:"user_name" validate:"required,alphanum"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"omitempty,min=2,alphanum"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,alphanum"`
	Hash      string `json:"hash" validate:"required,base64"`
}

type LoginInput struct {
	ClientID  string `json:"client_id" validate:"required"`
	UserName  string `json:"user_name" validate:"required,alphanum"`
	Hash      string `json:"hash" validate:"required,base64"`
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