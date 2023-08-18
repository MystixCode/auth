package user

import (
	"auth/conf"
	"auth/log"

	"time"
	"strconv"
	"io/ioutil"
	//"crypto/ed25519"
	//"encoding/pem"
	//"encoding/asn1"
	//"crypto/rsa"
	//"crypto/x509"

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

func (s *Service) Login(input LoginInput) (*TokenResponse, error) {
	//////////////////////////////////////////////////////////////////////////////////////////
	// Validate and verify input
	//////////////////////////////////////////////////////////////////////////////////////////
	err := s.Validator.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			s.Log.Error().Err(err).Msg("Validation failed")
		}
		return nil, ErrValidationFailed
	}

	// Verify user_name exists
	foundUser, err := s.Store.GetByUserName(input.UserName)
	if err != nil {
		s.Log.Error().Err(err).Msg("Verifying that user exists")
		return nil, ErrLogin
	}
	s.Log.Debug().Msgf("Found user: ", foundUser.UserName)

	// Verify hash
	if input.Hash != foundUser.Hash {
		s.Log.Error().Err(ErrLogin).Msg("Hashes are not equal")
		return nil, ErrLogin
	}
	s.Log.Debug().Msg("hash is equal")

	//////////////////////////////////////////////////////////////////////////////////////////
	// Get and parse signKey --> todo: create a function for that
	//////////////////////////////////////////////////////////////////////////////////////////

	var (

		signMethod  jwt.SigningMethod
		privPath	string
		signBytes   []byte
		signKey     interface{} // This will hold the private key or secret
	)

	//TODO: get AppName from app the user belongs to
	appName := "test"
	//TODO: get method from app the user belongs to
	signMethod = jwt.SigningMethodRS256


	// Determine the signing method
	if signMethod == jwt.SigningMethodRS256 {
		privPath = appName + "_rsa.pem"
		signBytes, err = ioutil.ReadFile(privPath)
		if err != nil {
			s.Log.Fatal().Err(err).Msg("Error reading private key bytes")
			return nil, err
		}
		parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
		if err != nil {
			s.Log.Fatal().Err(err).Msg("Error parsing private key")
			return nil, err
		}
		signKey = parsedKey
		s.Log.Debug().Msgf("PrivateKey RSA: %s", signKey)
	} else if signMethod == jwt.SigningMethodEdDSA {
		privPath = appName + "_ed25519.pem"
		signBytes, err = ioutil.ReadFile(privPath)
		if err != nil {
			s.Log.Fatal().Err(err).Msg("Error reading private key bytes")
			return nil, err
		}
		parsedKey, err := jwt.ParseEdPrivateKeyFromPEM(signBytes)
		if err != nil {
			s.Log.Fatal().Err(err).Msg("Error parsing private key")
			return nil, err
		}
		signKey = parsedKey
		s.Log.Debug().Msgf("PrivateKey Ed25519: %s", signKey)
	} else if signMethod == jwt.SigningMethodHS256 {
		signKey = []byte("AllYourBase")
	}


	//////////////////////////////////////////////////////////////////////////////////////////
	// Generate Access Token --> todo: create a function for that
	//////////////////////////////////////////////////////////////////////////////////////////

	accessTokenClaims := MyCustomClaims{
		foundUser.UserName,
		"profileRead",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "localhost:8080/v1",
			Subject:   strconv.Itoa(foundUser.ID),
			ID:        uuid.New().String(),
			Audience:  []string{"game_server"},
		},
	}
	accessToken := jwt.NewWithClaims(signMethod, accessTokenClaims)
	s.Log.Debug().Msgf("Access token generated: %v", accessToken)
	ssAccess, err := accessToken.SignedString(signKey)
	if err != nil {
		return nil, err
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	// Generate Refresh Token --> todo: create a function for that
	//////////////////////////////////////////////////////////////////////////////////////////

	refreshTokenClaims := jwt.RegisteredClaims{
		Issuer:    "https://example.com",  // Replace with your issuer URL
		Subject:   strconv.Itoa(foundUser.ID),         // Replace with the user ID
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // Refresh token expiration (30 days)
	}

	// Create a new token object with the specified claims
	refreshToken := jwt.NewWithClaims(signMethod, refreshTokenClaims)
	s.Log.Debug().Msgf("Refresh token generated: %v", refreshToken)
	ssRefresh, err := refreshToken.SignedString(signKey)
	if err != nil {
		return nil, err
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	// Create the response map
	//////////////////////////////////////////////////////////////////////////////////////////
	tokenResponse := &TokenResponse{
		TokenType:    "bearer",
		AccessToken:  ssAccess,
		RefreshToken: ssRefresh,
	}

	return tokenResponse, nil
}

func (s *Service) Create(input UserInput) (*User, error) {
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
