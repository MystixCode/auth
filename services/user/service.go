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


func (s *Service) readKey(clientID string, alg string) (signKey interface{}, err error) {
	s.Log.Debug().Msg("ReadKey Func ;)")
	//TODO: 5: this func in key service!!!

	var (
		privPath   string
		signBytes  []byte
	)

	switch alg {
	case "RS256":
		privPath = clientID + "_rsa.pem"
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
	case "Ed25519":
		privPath = clientID + "_ed25519.pem"
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
	case "HS256":
		signKey = []byte("AllYourBase")
	default:
		s.Log.Error().Err(ErrLogin).Msg("unknown alg")
	}

	return signKey, nil
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
	foundUser, err := s.Store.GetByUserName(input.UserName)
	if err != nil {
		s.Log.Error().Err(err).Msg("Verifying that user exists")
		return nil, ErrLogin
	}
	s.Log.Debug().Msgf("Found user: ", foundUser.UserName)


	// TODO: 2: verify client_id exists !!!

	//include appService
	//get app from appService.GetByClientID(input.ClientID)
	//then use foundApp for later request


	// TODO: 3: verify that the user belongs to app with ClientID !!!

	// Verify hash
	if input.Hash != foundUser.Hash {
		s.Log.Error().Err(ErrLogin).Msg("Hashes are not equal")
		return nil, ErrLogin
	}
	s.Log.Debug().Msg("hash is equal")

	//TODO: 4: get alg from app!!!

	alg := "Ed25519"

	//readKey
	signKey, err := s.readKey(input.ClientID, alg)
	if err != nil {
		s.Log.Fatal().Err(err).Msg("Error reading key")
		return nil, err
	}

	//TODO: get scopes for this user
	scopes := "readProfile"
	userID := foundUser.ID
	userName := foundUser.UserName

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
