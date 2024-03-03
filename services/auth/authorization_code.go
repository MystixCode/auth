package auth

import (
	// "auth/conf"
	// "auth/log"
	// "auth/services/key"

	"time"

	// "gorm.io/gorm"
	"github.com/go-playground/validator/v10"
)

type AuthorizationCode struct {
	ID        uint   `gorm:"primaryKey"`
	Code      string `json:"code"`
	Expiry    int64  `json:"expiry"`
	CreatedAt int64  `json:"created_at"`

	UserID    uint `json:"user_id" gorm:"foreignKey:UserID;references:User(ID)"`
	AppID     uint `json:"app_id" gorm:"foreignKey:AppID;references:App(ID);onDelete:CASCADE"` // CASCADE delete behavior
}

type AuthorizationCodeInput struct {
	Code  		string	`json:"code" validate:"required"`
	UserID		uint 	`json:"user_id" validate:"required"`
	AppID		uint 	`json:"app_id" validate:"required"`
}

func (s *Service) CreateAuthorizationCode(input AuthorizationCodeInput) (*AuthorizationCode, error) {
	var a AuthorizationCode

	// Validate input
	err := s.Validator.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			s.Log.Error().Err(err).Msg("Validation failed")
		}
		return nil, ErrValidationFailed
	}

	a.Code = input.Code
	a.Expiry = time.Now().Add(5 * time.Minute).Unix()
	a.CreatedAt = time.Now().Unix()
	a.UserID = input.UserID
	a.AppID = input.AppID
	//TODO add scopes that this athorizationCode authorizes? with foreign key !!!

	createdAuthorizationCode, err := s.AuthStore.CreateAuthorizationCode(&a)
	if err != nil {
		return nil, err
	}

	return createdAuthorizationCode, nil
}

func (s *Service) GetAuthorizationCodeByID(id string) (*AuthorizationCode, error) {
	return s.AuthStore.GetAuthorizationCodeByID(id)
}

func (s *Service) GetAllAuthorizationCodes() ([]*AuthorizationCode, error) {
	return s.AuthStore.GetAllAuthorizationCodes()
}

func (s *Service) DeleteAuthorizationCode(id string) error {
	return s.AuthStore.DeleteAuthorizationCode(id)
}
