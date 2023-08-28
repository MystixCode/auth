package key

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
	return &Store{
		log:  log,
		// conf: conf,
		db:   db,
	}
}

func (s *Store) Create(key *Key) (*Key, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	inserted := s.db.WithContext(ctx).Create(key)
	if inserted.Error != nil {
		return nil, ErrInsertFailed
	}

	return key, nil
}

func (s *Store) GetAll() ([]*Key, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var keys []*Key

	result := s.db.WithContext(ctx).Find(&keys)
	if result.Error != nil {
		return nil, ErrFindFailed
	}

	return keys, nil
}

func (s *Store) GetByID(id string) (*Key, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var key *Key

	result := s.db.WithContext(ctx).First(&key, id)
	if result.Error != nil {
		return nil, ErrNotFound
	}

	return key, nil
}

func (s *Store) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var key *Key

	result := s.db.WithContext(ctx).Delete(&key, id)
	if result.Error != nil {
		return ErrDeleteFailed
	} else if result.RowsAffected < 1 {
		return ErrNotFound

	}

	return nil
}
