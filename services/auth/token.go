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

//code=4%2F0Adeu5BXZmJQzaVv3EqyDxJ2uS9735DencVVvDL6aBG_illRGu5nbOmU77kA28u5913GcLQ
//&redirect_uri=https%3A%2F%2Fdevelopers.google.com%2Foauthplayground
//&client_id=407408718192.apps.googleusercontent.com
//&scope=
//&grant_type=authorization_code
type TokenInput struct {
	ClientID    string	`json:"client_id"`
	Code  		string `json:"code" validate:"required"`
	RedirectURI	string `json:"redirect_uri" validate:"required"`
	Scope  		string `json:"scope" validate:"required"`
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MyCustomClaims struct {
	Username string `json:"username"`
	Scopes string `json:"scopes"`
	jwt.RegisteredClaims
}

func (s *Service) generateAccessToken(userName string, userID uint, scopes string, alg string, signKey interface{}) (accessToken string, err error) {
	s.Log.Debug().Msg("generateAccessToken Func ;)")

	claims := MyCustomClaims{
		userName,
		scopes,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "localhost:8080/v1",
			Subject:   strconv.Itoa(int(userID)),
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
	//TODO: add to db !!!
	return accessToken, nil
}

func (s *Service) generateRefreshToken(userID uint, alg string, signKey interface{}) (refreshToken string, err error) {
	s.Log.Debug().Msg("generateRfreshToken Func ;)")

	claims := jwt.RegisteredClaims{
		Issuer:    "https://auth.example.com",  // Replace with your issuer URL
		Subject:   strconv.Itoa(int(userID)),         // Replace with the user ID
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
	//TODO: add to db !!!
	return refreshToken, nil
}	

func (s *Service) Token(input TokenInput) (*TokenResponse, error) {

	// Validate input
	err := s.Validator.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			s.Log.Error().Err(err).Msg("Validation failed")
		}
		return nil, ErrValidationFailed
	}

	//TODO: check authorizationcode exists and isnt expired !!! foundAuthorizationCode.Expiry > now
	foundAuthorizationCode, err := s.AuthStore.GetAuthorizationCodeByCode(input.Code)
	if err != nil {
		s.Log.Error().Err(err).Msg("Verifying that AuthorizationCode exists")
		return nil, ErrLogin
	}
	s.Log.Debug().Msgf("Found AuthorizationCode: ", foundAuthorizationCode.Code)

	if foundAuthorizationCode.Expiry <= time.Now().Unix() {
		s.Log.Error().Msg("AutorizationCode is expired")
		return nil, ErrLogin
	}

	// verify client_id/app exists
	foundApp, err := s.AuthStore.GetAppByClientID(input.ClientID)
	if err != nil {
		s.Log.Error().Err(err).Msg("Verifying that client_id exists")
		return nil, ErrLogin
	}
	s.Log.Debug().Msgf("Found client_id: ", foundApp.ClientID)

	// Verify and get user
	foundUser, err := s.AuthStore.GetUserByID(strconv.Itoa(int(foundAuthorizationCode.UserID)))
	if err != nil {
		s.Log.Error().Err(err).Msg("Verifying that user exists")
		return nil, ErrLogin
	}
	s.Log.Debug().Msgf("Found user: ", foundUser.UserName)

	// TODO: verify that the user belongs to app with ClientID !!!

	//TODO: check the user has permission for requested scopes !!!

	//TODO:
	//scopes := foundAuthorizationCode.Scopes !!!
	scopes := "testscopetodo"
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
