package auth

import (
// 	"auth/conf"
// 	"auth/log"

	"time"
	"strconv"

// 	"gorm.io/gorm"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        int    `json:"id" gorm:"primaryKey"`
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
	Hash      string `json:"hash" validate:"required,base64"`
}

type LoginInput struct {
	ClientID  string `json:"client_id" validate:"required"`
	UserName  string `json:"user_name" validate:"required,alphanum"`
	Hash      string `json:"hash" validate:"required,base64"`
}

func (s *Service) generateAccessToken(userName string, userID int, scopes string, alg string, signKey interface{}) (accessToken string, err error) {
	s.Log.Debug().Msg("generateAccessToken Func ;)")

	claims := MyCustomClaims{
		userName,
		scopes,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "localhost:8080/v1",
			Subject:   strconv.Itoa(userID),
			ID:        uuid.New().String(),
			Audience:  []string{"game_server"},
		},
	}

	var signMethod jwt.SigningMethod
	switch alg {
	case "RS256":
		signMethod = jwt.SigningMethodRS256
	case "Ed25519":
		signMethod = jwt.SigningMethodEdDSA
	case "HS256":
		signMethod = jwt.SigningMethodHS256
	default:
		s.Log.Error().Err(ErrLogin).Msg("unknown alg")
	}

	unsignedToken := jwt.NewWithClaims(signMethod, claims)
	
	accessToken, err = unsignedToken.SignedString(signKey)
	if err != nil {
		return "", err
	}

	s.Log.Debug().Msgf("Access token generated: %v", unsignedToken)
	return accessToken, nil
}

func (s *Service) generateRefreshToken(userID int, alg string, signKey interface{}) (refreshToken string, err error) {
	s.Log.Debug().Msg("generateRfreshToken Func ;)")

	claims := jwt.RegisteredClaims{
		Issuer:    "https://auth.example.com",  // Replace with your issuer URL
		Subject:   strconv.Itoa(userID),         // Replace with the user ID
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // Refresh token expiration (30 days)
	}

	var signMethod jwt.SigningMethod
	switch alg {
	case "RS256":
		signMethod = jwt.SigningMethodRS256
	case "Ed25519":
		signMethod = jwt.SigningMethodEdDSA
	case "HS256":
		signMethod = jwt.SigningMethodHS256
	default:
		s.Log.Error().Err(ErrLogin).Msg("unknown alg")
	}

	unsignedToken := jwt.NewWithClaims(signMethod, claims)
	
	refreshToken, err = unsignedToken.SignedString(signKey)
	if err != nil {
		return "", err
	}

	s.Log.Debug().Msgf("Refresh token generated: %v", unsignedToken)
	return refreshToken, nil
}	

func (s *Service) Login(input LoginInput) (*TokenResponse, error) {

	// Validate input
	err := s.Validator.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			s.Log.Error().Err(err).Msg("Validation failed")
		}
		return nil, ErrValidationFailed
	}

	// Verify user_name exists
	foundUser, err := s.AuthStore.GetUserByUserName(input.UserName)
	if err != nil {
		s.Log.Error().Err(err).Msg("Verifying that user exists")
		return nil, ErrLogin
	}
	s.Log.Debug().Msgf("Found user: ", foundUser.UserName)

	// verify client_id/app exists
	foundApp, err := s.AuthStore.GetAppByClientID(input.ClientID)
	if err != nil {
		s.Log.Error().Err(err).Msg("Verifying that client_id exists")
		return nil, ErrLogin
	}
	s.Log.Debug().Msgf("Found client_id: ", foundApp.ClientID)

	// TODO: 3: verify that the user belongs to app with ClientID !!!

	// Verify hash
	if input.Hash != foundUser.Hash {
		s.Log.Error().Err(ErrLogin).Msg("Hashes are not equal")
		return nil, ErrLogin
	}
	s.Log.Debug().Msg("hash is equal")

	//TODO: get scopes for this user !!!
	scopes := "readProfile"
	userID := foundUser.ID
	userName := foundUser.UserName
	alg := foundApp.Alg
	clientID := foundApp.ClientID
	
	//readKey
	signKey, err := s.readKey(clientID, alg)
	if err != nil {
		s.Log.Fatal().Err(err).Msg("Error reading key")
		return nil, err
	}

	//generateAccessToken
	accessToken, err := s.generateAccessToken(userName, userID, scopes, alg, signKey)
	if err != nil {
		s.Log.Fatal().Err(err).Msg("Error generating access token")
		return nil, err
	}

	//generateRefreshToken
	refreshToken, err := s.generateRefreshToken(userID, alg, signKey)
	if err != nil {
		s.Log.Fatal().Err(err).Msg("Error generating refresh token")
		return nil, err
	}

	tokenResponse := &TokenResponse{
		TokenType:    "bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokenResponse, nil
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
	u.Hash = input.Hash
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

	if input.Hash != "" {
		return nil, ErrPasswordChangeNotAllowed
	}

	u.UpdatedAt = time.Now().Unix()

	return s.AuthStore.UpdateUser(id, u)
}

func (s *Service) DeleteUser(id string) error {
	return s.AuthStore.DeleteUser(id)
}
