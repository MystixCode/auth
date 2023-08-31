package auth

import (
	// "auth/conf"
	// "auth/log"
	// "auth/services/key"

	"time"

	// "gorm.io/gorm"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

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

func (s *Service) CreateApp(input AppInput) (*App, error) {
	var a App

	// Validate input
	err := s.Validator.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			s.Log.Error().Err(err).Msg("Validation failed")
		}
		return nil, ErrValidationFailed
	}

	a.AppName = input.AppName
	a.AppURI = input.AppURI
	a.Alg = input.Alg

	if input.RedirectURI != "" {
		a.RedirectURI = input.RedirectURI
	}

	if input.ClientType != "" {
		a.ClientType = input.ClientType
	}

	// Generate client id
	a.ClientID = uuid.New().String()

	timeNow := time.Now().Unix()
	a.CreatedAt = timeNow
	a.UpdatedAt = timeNow

	createdApp, err := s.AuthStore.CreateApp(&a)
	if err != nil {
		return nil, err
	}

    keyInput := KeyInput{
        AppID:  createdApp.ID, // Use the actual app ID here
		Alg:	a.Alg,
    }
    createdKey, err := s.CreateKey(keyInput)
	if err != nil {
		return nil, err
	}
	s.Log.Debug().Msgf("Created Key", createdKey)

	return createdApp, nil
}

func (s *Service) GetAppByID(id string) (*App, error) {
	return s.AuthStore.GetAppByID(id)
}

func (s *Service) GetAllApps() ([]*App, error) {
	return s.AuthStore.GetAllApps()
}

func (s *Service) UpdateApp(id string, input *AppInput) (*App, error) {
	a, err := s.AuthStore.GetAppByID(id)
	if err != nil {
		return nil, err
	}

	// Validate input
	err = s.Validator.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			s.Log.Error().Err(err).Msg("Validation failed")
		}
		return nil, ErrValidationFailed
	}

	a.AppName = input.AppName

	if input.AppURI != "" {
		a.AppURI = input.AppURI
	}

	if input.RedirectURI != "" {
		a.RedirectURI = input.RedirectURI
	}

	if input.ClientType != "" {
		a.ClientType = input.ClientType
	}

	a.UpdatedAt = time.Now().Unix()

	return s.AuthStore.UpdateApp(id, a)
}

func (s *Service) DeleteApp(id string) error {
	return s.AuthStore.DeleteApp(id)
}