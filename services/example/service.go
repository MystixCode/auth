package example

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

func (s *Service) Create(input ExampleInput) (*Example, error) {
	var t Example

	//TODO validation
	if input.ExampleName != "" {
		t.ExampleName = input.ExampleName
	}

	if input.ExampleValue != "" {
		t.ExampleValue = input.ExampleValue
	}

	timeNow := time.Now().Unix()
	t.CreatedAt = timeNow
	t.UpdatedAt = timeNow

	createdExample, err := s.Store.Create(&t)
	if err != nil {
		return nil, err
	}

	return createdExample, nil
}

func (s *Service) GetByID(id string) (*Example, error) {
	return s.Store.GetByID(id)
}

func (s *Service) GetAll() ([]*Example, error) {
	return s.Store.GetAll()
}

func (s *Service) Update(id string, input *ExampleInput) (*Example, error) {
	t, err := s.Store.GetByID(id)
	if err != nil {
		return nil, err
	}

	//TODO validation
	if input.ExampleName != "" {
		t.ExampleName = input.ExampleName
	}

	if input.ExampleValue != "" {
		t.ExampleValue = input.ExampleValue
	}

	t.UpdatedAt = time.Now().Unix()

	return s.Store.Update(id, t)
}

func (s *Service) Delete(id string) error {
	return s.Store.Delete(id)
}
