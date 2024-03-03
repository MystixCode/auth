package auth

import (
// 	"auth/conf"
// 	"auth/log"
	"os"
	"time"
	"strconv"
	"crypto/ed25519"
	"io/ioutil"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

// 	"gorm.io/gorm"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type Key struct {
	ID         	uint	`json:"id" gorm:"primaryKey"`
	AppID		uint		`json:"app_id" gorm:"not null"`
	Alg			string	`json:"alg" gorm:"not null"`
	CreatedAt	int64	`json:"created_at" gorm:"not null"`
}

type KeyInput struct {
	AppID 	uint		`json:"app_id" validate:"required,number"`
	Alg		string	`json:"alg" validate:"required,alphanum"`
}

func (s *Service) CreateKey(input KeyInput) (*Key, error) {
	var k Key

	// Validate input
	err := s.Validator.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			s.Log.Error().Err(err).Msg("Validation failed")
		}
		return nil, ErrValidationFailed
	}

	k.AppID = input.AppID
	k.Alg = input.Alg
	timeNow := time.Now().Unix()
	k.CreatedAt = timeNow

	createdKey, err := s.AuthStore.CreateKey(&k)
	if err != nil {
		return nil, err
	}

	foundApp, err := s.GetAppByID(strconv.Itoa(int(k.AppID))) // Call key service's Create method
	if err != nil {
		return nil, err
	}
	s.Log.Debug().Msgf("Found ClientID", foundApp.ClientID)

	// https://golang-jwt.github.io/jwt/usage/signing_methods/

	// Generate the keys for the app based on algorithm
	s.generateKeys(foundApp.ClientID,k.Alg)

	return createdKey, nil
}

func (s *Service) GetKeyByID(id string) (*Key, error) {
	return s.AuthStore.GetKeyByID(id)
}

func (s *Service) GetAllKeys() ([]*Key, error) {
	return s.AuthStore.GetAllKeys()
}

func (s *Service) DeleteKey(id string) error {
	return s.AuthStore.DeleteKey(id)
}


func (s *Service) generateRSAKeys(privPath string, pubPath string) (error) {

	var (
		err   error
		size  int
		b     []byte
		block *pem.Block
		pub   *rsa.PublicKey
		priv  *rsa.PrivateKey
	)

	size = 2048 // Replace with your desired key size or pass var to function...

	priv, err = rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return err
	}

	b = x509.MarshalPKCS1PrivateKey(priv)

	block = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: b,
	}

	err = ioutil.WriteFile(privPath, pem.EncodeToMemory(block), 0600)
	if err != nil {
		return err
	}

	// public key
	pub = &priv.PublicKey
	b, err = x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return err
	}

	block = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: b,
	}

	err = ioutil.WriteFile(pubPath, pem.EncodeToMemory(block), 0644)

	return err
}

func (s *Service) generateEd25519Keys(privPath string, pubPath string) error {

	var (
		err   error
		b     []byte
		block *pem.Block
		pub   ed25519.PublicKey
		priv  ed25519.PrivateKey
	)

	pub, priv, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	b, err = x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}

	block = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	}

	err = ioutil.WriteFile(privPath, pem.EncodeToMemory(block), 0600)
	if err != nil {
		return err
	}

	// public key
	b, err = x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return err
	}

	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}

	err = ioutil.WriteFile(pubPath, pem.EncodeToMemory(block), 0644)

	return err
}

func (s *Service) generateHmacKey(privPath string) error {

	//TODO: generate a proper hmac key !!!
	secret := "te3st!?fdg123_hfjghk!g78jtest_CHANGE_THAT"
	message := "me0s-sage?mes566sage__CHANGE_THAT"

	b := []byte(secret)
	h := hmac.New(sha256.New, b)
	h.Write([]byte(message))
	signKey :=  []byte(base64.StdEncoding.EncodeToString(h.Sum(nil)))

	//signKey := []byte("AllYourBase;)")
	err := os.WriteFile(privPath, signKey, 0600)

	return err
}

func (s *Service) generateKeys(clientID string, alg string) error {
	var err error
	var keyDir string = "keys/"
	switch alg {
	case "RS256":
		privPath := keyDir + clientID + "_RS256.pem"
		pubPath := keyDir + clientID + "_RS256.pub.pem"
		err = s.generateRSAKeys(privPath, pubPath)
	case "Ed25519":
		privPath := keyDir + clientID + "_ed25519.pem"
		pubPath := keyDir + clientID + "_ed25519.pub.pem"
		err = s.generateEd25519Keys(privPath, pubPath)
	case "HS256":
		privPath := keyDir + clientID + "_HS256.key"
		err = s.generateHmacKey(privPath)
	default:
		s.Log.Error().Err(ErrKeyGenFailed).Msg("unknown alg")
	}

	return err
}

func (s *Service) readKey(clientID string, alg string) (signKey interface{}, err error) {
	s.Log.Debug().Msg("ReadKey Func ;)")

	var (
		keyDir      string = "keys/"
		privPath    string
		signBytes   []byte
	)

	switch alg {
	case "RS256":
		privPath = keyDir + clientID + "_RS256.pem"
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
		privPath = keyDir + clientID + "_ed25519.pem"
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
		privPath = keyDir + clientID + "_HS256.key"
		signKey, err = ioutil.ReadFile(privPath)
		if err != nil {
			s.Log.Fatal().Err(err).Msg("Error reading private key bytes")
			return nil, err
		}
		s.Log.Debug().Msgf("PrivateKey HS256: %s", signKey)

	default:
		s.Log.Error().Err(ErrLogin).Msg("unknown alg")
	}

	return signKey, nil
}
