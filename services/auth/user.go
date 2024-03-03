package auth

import (
// 	"auth/conf"
// 	"auth/log"

	"time"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserName  string `json:"user_name" gorm:"not null"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Hash      string `json:"hash" gorm:"not null"`
	CreatedAt int64  `json:"created_at" gorm:"not null"`
	UpdatedAt int64  `json:"updated_at"`
}

type UserInput struct {
	UserName  string `json:"user_name" validate:"required,alphanum"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"omitempty,min=2,alphanum"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,alphanum"`
	Password  string `json:"password" validate:"required,alphanum"`
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func (s *Service) CreateUser(input UserInput) (*User, error) {
	var u User

	//validation
	// returns nil or ValidationErrors ( []FieldError )
	err := s.Validator.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			s.Log.Error().Err(err).Msg("Validation failed")
		}
		return nil, ErrValidationFailed
	}

	u.UserName = input.UserName
	u.Email = input.Email
	u.FirstName = input.FirstName
	u.LastName = input.LastName
	hash, err := HashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	u.Hash = hash
	timeNow := time.Now().Unix()
	u.CreatedAt = timeNow
	u.UpdatedAt = timeNow

	createdUser, err := s.AuthStore.CreateUser(&u)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *Service) GetUserByID(id string) (*User, error) {
	return s.AuthStore.GetUserByID(id)
}

func (s *Service) GetAllUsers() ([]*User, error) {
	return s.AuthStore.GetAllUsers()
}

func (s *Service) UpdateUser(id string, input *UserInput) (*User, error) {
	u, err := s.AuthStore.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	err = s.Validator.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			s.Log.Error().Err(err).Msg("Validation failed")
		}
		return nil, ErrValidationFailed
	}

	if input.UserName != "" {
		u.UserName = input.UserName
	}

	if input.Email != "" {
		u.Email = input.Email
	}

	if input.FirstName != "" {
		u.FirstName = input.FirstName
	}

	if input.LastName != "" {
		u.LastName = input.LastName
	}

	if input.Password != "" {
		return nil, ErrPasswordChangeNotAllowed
	}

	u.UpdatedAt = time.Now().Unix()

	return s.AuthStore.UpdateUser(id, u)
}

func (s *Service) DeleteUser(id string) error {
	return s.AuthStore.DeleteUser(id)
}
