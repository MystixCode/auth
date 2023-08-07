package authorize

import (
	"auth/conf"
	"auth/log"

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

func (s *Store) Create(authorize *AuthorizationCode) (*AuthorizationCode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	inserted := s.db.WithContext(ctx).Create(authorize)
	if inserted.Error != nil {

		return nil, ErrInsertFailed
	}

	return authorize, nil
}

func (s *Store) GetAll() ([]*AuthorizationCode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var authorizes []*AuthorizationCode

	result := s.db.WithContext(ctx).Find(&authorizes)
	if result.Error != nil {
		return nil, ErrFindFailed
	}

	return authorizes, nil
}

func (s *Store) GetByID(id string) (*AuthorizationCode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var authorize *AuthorizationCode

	result := s.db.WithContext(ctx).First(&authorize, id)
	if result.Error != nil {
		return nil, ErrNotFound
	}

	return authorize, nil
}

// func (s *Store) GetByUsername(username string) (*User, error) {
// 	return s.getByKeyValue("username", username)
// }

// func (s *Store) GetByEmail(email string) (*User, error) {
// 	return s.getByKeyValue("email", email)
// }

func (s *Store) Update(id string, authorize *AuthorizationCode) (*AuthorizationCode, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result := s.db.WithContext(ctx).Save(&authorize)
	if result.Error != nil {
		return nil, ErrUpdatedFailed
	}

	return authorize, nil
}

func (s *Store) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var authorize *AuthorizationCode

	result := s.db.WithContext(ctx).Delete(&authorize, id)
	if result.Error != nil {
		return ErrDeleteFailed
	} else if result.RowsAffected < 1 {
		return ErrNotFound

	}

	return nil
}

// func (s *Store) getByKeyValue(key string, value interface{}) (*User, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 	defer cancel()

// 	var user User

// 	err := s.collection.FindOne(ctx, bson.M{key: value}).Decode(&user)
// 	if err != nil {
// 		switch err {
// 		case mongo.ErrNoDocuments:
// 			return nil, ErrNotFound
// 		default:
// 			s.logger.Warnf("error while finding user: %v", err)
// 			return nil, ErrFindFailed
// 		}
// 	}

// 	return &user, nil
// }
