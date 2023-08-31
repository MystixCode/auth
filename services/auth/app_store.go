package auth

import (
	// "auth/conf"
	// "auth/log"

	"context"
	"time"

	// "gorm.io/gorm"
)


func (s *AuthStore) CreateApp(app *App) (*App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	inserted := s.db.WithContext(ctx).Create(app)
	if inserted.Error != nil {

		return nil, ErrInsertFailed
	}

	return app, nil
}

func (s *AuthStore) GetAllApps() ([]*App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var apps []*App

	result := s.db.WithContext(ctx).Find(&apps)
	if result.Error != nil {
		return nil, ErrFindFailed
	}

	return apps, nil
}

func (s *AuthStore) GetAppByID(id string) (*App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var app *App

	result := s.db.WithContext(ctx).First(&app, id)
	if result.Error != nil {
		return nil, ErrNotFound
	}

	return app, nil
}

// func (s *Store) GetByUsername(username string) (*User, error) {
// 	return s.getByKeyValue("username", username)
// }

// func (s *Store) GetByEmail(email string) (*User, error) {
// 	return s.getByKeyValue("email", email)
// }

func (s *AuthStore) UpdateApp(id string, app *App) (*App, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result := s.db.WithContext(ctx).Save(&app)
	if result.Error != nil {
		return nil, ErrUpdatedFailed
	}

	return app, nil
}

func (s *AuthStore) DeleteApp(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var app *App

	result := s.db.WithContext(ctx).Delete(&app, id)
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
