package app

import (
	"auth/conf"
	"auth/log"

	"time"

	"gorm.io/gorm"
)

type Service struct {
	Log   *log.Logger
	Store *Store
}

func NewService(log *log.Logger, conf *conf.Config, db *gorm.DB) *Service {
	return &Service{
		Log:   log,
		Store: NewStore(log, conf, db),
	}
}

func (s *Service) Create(input AppInput) (*App, error) {
	var a App

	//TODO validation
	if input.AppName != "" {
		a.AppName = input.AppName
	}

	if input.AppURL != "" {
		a.AppURL = input.AppURL
	}

	if input.RedirectURL != "" {
		a.RedirectURL = input.RedirectURL
	}

	if input.ClientType != "" {
		a.ClientType = input.ClientType
	}

	timeNow := time.Now().Unix()
	a.CreatedAt = timeNow
	a.UpdatedAt = timeNow

	createdApp, err := s.Store.Create(&a)
	if err != nil {
		return nil, err
	}

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

	if input.AppName != "" {
		a.AppName = input.AppName
	}

	if input.AppURL != "" {
		a.AppURL = input.AppURL
	}

	if input.RedirectURL != "" {
		a.RedirectURL = input.RedirectURL
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
