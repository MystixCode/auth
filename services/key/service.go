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
)

type Service struct {
	Log   *log.Logger
	Store *Store
}

func NewService(log *log.Logger, conf *conf.Config, db *gorm.DB) *Service {
	return &Service{
		Log:   log,
		Store: NewStore(log, conf, db),
	}
}

func (s *Service) Create(input KeyInput) (*Key, error) {
	var k Key

	//TODO validation

	if input.AppID != 0 {
		k.AppID = input.AppID
	}

	timeNow := time.Now().Unix()
	k.CreatedAt = timeNow

	createdKey, err := s.Store.Create(&k)
	if err != nil {
		return nil, err
	}

	//TODO get client id from App
	var clientID = 666
	var signMethod = "EdDSA"
	s.Generate(clientID,signMethod)

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

func (s *Service) Generate(clientID int, method string) error {
	var err error

	switch method {
	case "RS256":
		privPath := strconv.Itoa(clientID) + "_rsa.pem"
		pubPath := strconv.Itoa(clientID) + "_rsa.pub.pem"
		err = generateRSAKeys(privPath, pubPath)

	case "EdDSA":
		privPath := strconv.Itoa(clientID) + "_ed25519.pem"
		pubPath := strconv.Itoa(clientID) + "_ed25519.pub.pem"
		err = generateEd25519Keys(privPath, pubPath)

	case "HS256":
		s.Log.Debug().Msg("todo: function to generate and save to file")

	default:
		s.Log.Debug().Msg("unknown method")
	}

	return err
}
