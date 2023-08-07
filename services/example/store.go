package example

import (
	"auth/log"
	"auth/conf"

	"context"
	"time"

	"gorm.io/gorm"
)

type Store struct {
	log  *log.Logger
	conf *conf.Config
	db   *gorm.DB
}

func NewStore(log *log.Logger, conf *conf.Config, db *gorm.DB) *Store {

	//TODO: is the conf really needed here?
	return &Store{
		log:  log,
		conf: conf,
		db:   db,
	}
}

func (s *Store) Create(example *Example) (*Example, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	inserted := s.db.WithContext(ctx).Create(example)
	if inserted.Error != nil {
		return nil, ErrInsertFailed
	}

	return example, nil
}

func (s *Store) GetAll() ([]*Example, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var examples []*Example

	result := s.db.WithContext(ctx).Find(&examples)
	if result.Error != nil {
		return nil, ErrFindFailed
	}

	return examples, nil
}

func (s *Store) GetByID(id string) (*Example, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var example *Example

	result := s.db.WithContext(ctx).First(&example, id)
	if result.Error != nil {
		return nil, ErrNotFound
	}

	return example, nil
}

func (s *Store) Update(id string, example *Example) (*Example, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result := s.db.WithContext(ctx).Save(&example)
	if result.Error != nil {
		return nil, ErrUpdatedFailed
	}

	return example, nil
}

func (s *Store) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var example *Example

	result := s.db.WithContext(ctx).Delete(&example, id)
	if result.Error != nil {
		return ErrDeleteFailed
	} else if result.RowsAffected < 1 {
		return ErrNotFound

	}

	return nil
}
