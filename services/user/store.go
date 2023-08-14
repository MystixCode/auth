package user

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

func (s *Store) Create(user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	inserted := s.db.WithContext(ctx).Create(user)
	if inserted.Error != nil {

		return nil, ErrInsertFailed
	}

	return user, nil
}

func (s *Store) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var users []*User

	result := s.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, ErrFindFailed
	}

	return users, nil
}

func (s *Store) GetByID(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user *User

	result := s.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		return nil, ErrNotFound
	}

	return user, nil
}

func (s *Store) GetByUserName(userName string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user *User

	result := s.db.WithContext(ctx).Where("user_name = ?", userName).First(&user)
	if result.Error != nil {
		return nil, ErrLogin
	}

	return user, nil
}

// func (s *Store) GetByEmail(email string) (*User, error) {
// 	return s.getByKeyValue("email", email)
// }

func (s *Store) Update(id string, user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result := s.db.WithContext(ctx).Save(&user)
	if result.Error != nil {
		return nil, ErrUpdatedFailed
	}

	return user, nil
}

func (s *Store) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user *User

	result := s.db.WithContext(ctx).Delete(&user, id)
	if result.Error != nil {
		return ErrDeleteFailed
	} else if result.RowsAffected < 1 {
		return ErrNotFound

	}

	return nil
}

// func (s *Store) getByKeyValue(key string, value interface{}) (*User, error) {
//  	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
//  	defer cancel()
//  	var user User

//  	err := s.collection.FindOne(ctx, bson.M{key: value}).Decode(&user)
//  	if err != nil {
//  		switch err {
//  		case mongo.ErrNoDocuments:
//  			return nil, ErrNotFound
//  		default:
//  			s.logger.Warnf("error while finding user: %v", err)
//  			return nil, ErrFindFailed
//  		}
//  	}

//  	return &user, nil
// }
