package app

import (
	"auth/conf"
	"auth/log"
	"auth/services/key"

	"time"

	"gorm.io/gorm"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Service struct {
	Log   *log.Logger
	Store *Store
	Validator *validator.Validate
	KeyService *key.Service // Inject key service dependency
}

func NewService(log *log.Logger, conf *conf.Config, db *gorm.DB, validator *validator.Validate, keyService *key.Service) *Service {
	return &Service{
		Log:   log,
		Store: NewStore(log, conf, db),
		Validator: validator,
		KeyService: keyService,
	}
}

func (s *Service) Create(input AppInput) (*App, error) {
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

	// if input.RedirectURL != "" {
	// 	a.RedirectURL = input.RedirectURL
	// }

	// if input.ClientType != "" {
	// 	a.ClientType = input.ClientType
	// }

	// Generate client id
	a.ClientID = uuid.New().String()
	timeNow := time.Now().Unix()
	a.CreatedAt = timeNow
	a.UpdatedAt = timeNow

	createdApp, err := s.Store.Create(&a)
	if err != nil {
		return nil, err
	}

    // Create key using the injected key service
    keyInput := key.KeyInput{
        AppID:  createdApp.ID, // Use the actual app ID here
		Alg:	a.Alg,
    }
    createdKey, err := s.KeyService.Create(keyInput) // Call key service's Create method
	if err != nil {
		return nil, err
	}
	s.Log.Debug().Msgf("Created Key", createdKey)

	return createdApp, nil
}

func (s *Service) GetByID(id string) (*App, error) {

	return s.Store.GetByID(id)

}

func (s *Service) GetAll() ([]*App, error) {

	return s.Store.GetAll()

}

func (s *Service) Update(id string, input *AppInput) (*App, error) {
	a, err := s.Store.GetByID(id)
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

	return s.Store.Update(id, a)
}

func (s *Service) Delete(id string) error {

	return s.Store.Delete(id)
}

// func (s *Service) verifyPassword(hash string, password string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	if err != nil {
// 		// TODO add logging
// 		return false
// 	}

// 	return true
// }
