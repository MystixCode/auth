package auth

import (
	// "auth/conf"
	// "auth/log"

	"context"
	"time"

	// "gorm.io/gorm"
)


func (s *AuthStore) CreateAuthorizationCode(authorizationCode *AuthorizationCode) (*AuthorizationCode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	inserted := s.db.WithContext(ctx).Create(authorizationCode)
	if inserted.Error != nil {

		return nil, ErrInsertFailed
	}

	return authorizationCode, nil
}

func (s *AuthStore) GetAllAuthorizationCodes() ([]*AuthorizationCode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var authorizationCodes []*AuthorizationCode

	result := s.db.WithContext(ctx).Find(&authorizationCodes)
	if result.Error != nil {
		return nil, ErrFindFailed
	}

	return authorizationCodes, nil
}

func (s *AuthStore) GetAuthorizationCodeByID(id string) (*AuthorizationCode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var authorizationCode *AuthorizationCode

	result := s.db.WithContext(ctx).First(&authorizationCode, id)
	if result.Error != nil {
		return nil, ErrNotFound
	}

	return authorizationCode, nil
}

func (s *AuthStore) GetAuthorizationCodeByCode(code string) (*AuthorizationCode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var authorizationCode *AuthorizationCode

	result := s.db.WithContext(ctx).Where("code = ?", code).First(&authorizationCode)
	if result.Error != nil {
		return nil, ErrLogin
	}

	return authorizationCode, nil
}

func (s *AuthStore) DeleteAuthorizationCode(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var authorizationCode *AuthorizationCode

	result := s.db.WithContext(ctx).Delete(&authorizationCode, id)
	if result.Error != nil {
		return ErrDeleteFailed
	} else if result.RowsAffected < 1 {
		return ErrNotFound

	}

	return nil
}
