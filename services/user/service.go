package user

import (
	"auth/conf"
	"auth/log"

	"time"

	"fmt"



	"gorm.io/gorm"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
)


type Service struct {
	Log   *log.Logger
	Store *Store
	Validator *validator.Validate
}

func NewService(log *log.Logger, conf *conf.Config, db *gorm.DB, validator *validator.Validate) *Service {
	return &Service{
		Log:   log,
		Store: NewStore(log, conf, db),
		Validator: validator,
	}
}

type MyCustomClaims struct {
	Username string `json:"username"`
	Scopes string `json:"scopes"`
	jwt.RegisteredClaims
}

func (s *Service) Login(input LoginInput) (*TokenResponse, error) {
	var u LoginInput

	//TODO validation
	if input.Email != "" {
		u.Email = input.Email
	}

	if input.Hash != "" {
		u.Hash = input.Hash
	}


	// returns nil or ValidationErrors ( []FieldError )
	err := s.Validator.Struct(u)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return nil, err
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}

		// from here you can create your own error messages in whatever language you wish
		return nil, err
	}



	// TODO:

	// Verify email

	// Verify hash

	//////////////////////////////////////////////////////////////////////////////////////////
	// Generate Access Token
	//////////////////////////////////////////////////////////////////////////////////////////

	mySigningKey := []byte("AllYourBase")

	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		"Spartan117",
		"profileRead",
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "localhost:8080/v1",
			Subject:   "uid",
			ID:        uuid.New().String(),
			Audience:  []string{"game_server"},
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ssAccess, err := accessToken.SignedString(mySigningKey)
	if err != nil {
		return nil, err
	}
	//////////////////////////////////////////////////////////////////////////////////////////
	// Generate Refresh Token
	//////////////////////////////////////////////////////////////////////////////////////////

	// Create claims for the refresh token
	refreshTokenClaims := jwt.RegisteredClaims{
		// You can add any relevant claims here
		Issuer:    "https://example.com",  // Replace with your issuer URL
		Subject:   "user_id_here",         // Replace with the user ID
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // Refresh token expiration (30 days)
	}

	// Create a new token object with the specified claims
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	ssRefresh, err := refreshToken.SignedString(mySigningKey)
	if err != nil {
		return nil, err
	}
	//////////////////////////////////////////////////////////////////////////////////////////

	// Create the response map
	tokenResponse := &TokenResponse{
		TokenType:    "bearer",
		AccessToken:  ssAccess,
		RefreshToken: ssRefresh,
	}

	return tokenResponse, nil
}

func (s *Service) Create(input UserInput) (*User, error) {
	var u User

	//TODO validation
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

	if input.Hash != "" {
		u.Hash = input.Hash
	}

	timeNow := time.Now().Unix()
	u.CreatedAt = timeNow
	u.UpdatedAt = timeNow

	createdUser, err := s.Store.Create(&u)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *Service) GetByID(id string) (*User, error) {

	return s.Store.GetByID(id)

}

func (s *Service) GetAll() ([]*User, error) {

	return s.Store.GetAll()

}

func (s *Service) Update(id string, input *UserInput) (*User, error) {
	u, err := s.Store.GetByID(id)
	if err != nil {
		return nil, err
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

	if input.Hash != "" {
		return nil, ErrPasswordChangeNotAllowed
	}

	u.UpdatedAt = time.Now().Unix()

	return s.Store.Update(id, u)
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
