package app

import (
	"auth/conf"
	"auth/log"
	"auth/services/key"

	"time"
	"crypto/ed25519"
	"io/ioutil"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	//"os"
	//"path/filepath"

	"gorm.io/gorm"
)

type Service struct {
	Log   *log.Logger
	Store *Store
	KeyService *key.Service // Inject key service dependency
}

func NewService(log *log.Logger, conf *conf.Config, db *gorm.DB, keyService *key.Service) *Service {
	return &Service{
		Log:   log,
		Store: NewStore(log, conf, db),
		KeyService: keyService,
	}
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

func (s *Service) Create(input AppInput) (*App, error) {
	var a App

	//TODO validation
	if input.AppName != "" {
		a.AppName = input.AppName
	}

	if input.AppURI != "" {
		a.AppURI = input.AppURI
	}

	if input.SignMethod != "" {
		a.SignMethod = input.SignMethod
	}

	// if input.RedirectURL != "" {
	// 	a.RedirectURL = input.RedirectURL
	// }

	// if input.ClientType != "" {
	// 	a.ClientType = input.ClientType
	// }

	timeNow := time.Now().Unix()
	a.CreatedAt = timeNow
	a.UpdatedAt = timeNow

	createdApp, err := s.Store.Create(&a)
	if err != nil {
		return nil, err
	}

    // Create key using the injected key service
    keyInput := key.KeyInput{
        AppID:  createdApp.ID, // Use the actual app ID here
    }
    createdKey, err := s.KeyService.Create(keyInput) // Call key service's Create method
	if err != nil {
		return nil, err
	}
	s.Log.Debug().Msgf("Created Key", createdKey)

	return createdApp, nil
}

func (s *Service) GetByID(id string) (*App, error) {

	return s.Store.GetByID(id)

}

func (s *Service) GetAll() ([]*App, error) {

	return s.Store.GetAll()

}

func (s *Service) Update(id string, input *AppInput) (*App, error) {
	a, err := s.Store.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.AppName != "" {
		a.AppName = input.AppName
	}

	if input.AppURI != "" {
		a.AppURI = input.AppURI
	}

	// if input.RedirectURL != "" {
	// 	a.RedirectURL = input.RedirectURL
	// }

	// if input.ClientType != "" {
	// 	a.ClientType = input.ClientType
	// }

	a.UpdatedAt = time.Now().Unix()

	return s.Store.Update(id, a)
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
