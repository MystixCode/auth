package key

import (
	"auth/conf"
	"auth/log"

	"time"
	"strconv"
	"crypto/ed25519"
	"io/ioutil"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	//"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
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

func (s *Service) Create(input KeyInput) (*Key, error) {
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

	createdKey, err := s.Store.Create(&k)
	if err != nil {
		return nil, err
	}

	// TODO: 1: get client id !!!

	//include app service
	//get app from appService.GetByID(k.AppID)
	//then use result.ClientID

	// https://golang-jwt.github.io/jwt/usage/signing_methods/

	var clientID = 666

	// Generate the keys for the app based on algorithm
	s.generate(clientID,k.Alg)

	return createdKey, nil
}

func (s *Service) GetByID(id string) (*Key, error) {
	return s.Store.GetByID(id)
}

func (s *Service) GetAll() ([]*Key, error) {
	return s.Store.GetAll()
}

func (s *Service) Delete(id string) error {
	return s.Store.Delete(id)
}


func generateRSAKeys(privPath string, pubPath string) (error) {

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

func generateEd25519Keys(privPath string, pubPath string) error {

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

func (s *Service) generate(clientID int, alg string) error {
	var err error

	switch alg {
	case "RS256":
		privPath := strconv.Itoa(clientID) + "_rsa.pem"
		pubPath := strconv.Itoa(clientID) + "_rsa.pub.pem"
		err = generateRSAKeys(privPath, pubPath)
	case "Ed25519":
		privPath := strconv.Itoa(clientID) + "_ed25519.pem"
		pubPath := strconv.Itoa(clientID) + "_ed25519.pub.pem"
		err = generateEd25519Keys(privPath, pubPath)
	case "HS256":
		s.Log.Debug().Msg("TODO: 6: function to generate and save to file!!!")

	default:
		s.Log.Error().Err(ErrKeyGenFailed).Msg("unknown alg")
	}

	return err
}
