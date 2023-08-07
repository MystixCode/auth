package authorize

import (
	"auth/conf"
	"auth/log"

	"net/http"
	"time"
	// "fmt"

	"gorm.io/gorm"
)

var (
	userAuthenticated = false
	userConsent       = false
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

func (s *Service) ViaPost(input AuthorizeInput, w http.ResponseWriter, r *http.Request) (string, error) {
	var a AuthorizeInput
	// Parse authorization request parameters
	if input.ClientID != "" {
		a.ClientID = input.ClientID
	}
	if input.RedirectURL != "" {
		a.RedirectURL = input.RedirectURL
	}
	if input.ResponseType != "" {
		a.ResponseType = input.ResponseType
	}
	if input.Scope != "" {
		a.Scope = input.Scope
	}
	if input.State != "" {
		a.State = input.State
	}
	if input.CodeChallenge != "" {
		a.CodeChallenge = input.CodeChallenge
	}
	if input.CodeChallengeMethod != "" {
		a.CodeChallengeMethod = input.CodeChallengeMethod
	}
	//TODO validation

	//TODO !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	response := "Hello from authorize service.ClientID:" + string(a.ClientID) + " RedirectURL: " + string(a.RedirectURL) + " ResponseType: " + string(a.ResponseType) + " Scope: " + string(a.Scope) + " State: " + string(a.State) + " CodeChallenge: " + string(a.CodeChallenge + " CodeChallengeMethod: " + string(a.CodeChallengeMethod))


	// Check if user is authenticated and has given consent
	if !userAuthenticated || !userConsent {
		if !userAuthenticated {
			http.Redirect(w, r, "/v1/login", http.StatusFound)
		} else {
			http.Redirect(w, r, "/v1/consent", http.StatusFound)
		}
		return response, nil
	}

	// Generate authorization code
	//authorizationCode := "some-authorization-code"

	// Construct redirect URL with authorization code
	//string redirectURL
	//redirectURL := fmt.Sprintf("%s?code=%s&state=%s", RedirectURL, AuthorizationCode, State)
	var test = "https://localhost:3000/auth/callback/oder/so"

	// Perform the redirect
	http.Redirect(w, r, test, http.StatusFound)




	return response, nil
}












func (s *Service) ViaGet(input AuthorizeInput, w http.ResponseWriter, r *http.Request) (string, error) {
	var a AuthorizeInput


	if input.ClientID != "" {
		a.ClientID = input.ClientID
	}
	if input.RedirectURL != "" {
		a.RedirectURL = input.RedirectURL
	}
	if input.ResponseType != "" {
		a.ResponseType = input.ResponseType
	}
	if input.Scope != "" {
		a.Scope = input.Scope
	}
	if input.State != "" {
		a.State = input.State
	}
	if input.CodeChallenge != "" {
		a.CodeChallenge = input.CodeChallenge
	}
	if input.CodeChallengeMethod != "" {
		a.CodeChallengeMethod = input.CodeChallengeMethod
	}
	//TODO validation

	//TODO !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	response := "Hello from authorize service via get.ClientID:" + string(a.ClientID) + " RedirectURL: " + string(a.RedirectURL) + " ResponseType: " + string(a.ResponseType) + " Scope: " + string(a.Scope) + " State: " + string(a.State) + " CodeChallenge: " + string(a.CodeChallenge + " CodeChallengeMethod: " + string(a.CodeChallengeMethod))


	// return response, nil


	// Check if user is authenticated and has given consent
	if !userAuthenticated || !userConsent {
		if !userAuthenticated {
			http.Redirect(w, r, "/v1/login", http.StatusFound)
		} else {
			http.Redirect(w, r, "/v1/consent", http.StatusFound)
		}
		return response, nil
	}

	// Generate authorization code
	//authorizationCode := "some-authorization-code"

	// Construct redirect URL with authorization code
	//string redirectURL
	//redirectURL := fmt.Sprintf("%s?code=%s&state=%s", RedirectURL, AuthorizationCode, State)
	var test = "https://localhost:3000/auth/callback/oder/so"

	// Perform the redirect
	http.Redirect(w, r, test, http.StatusFound)




	return response, nil

}

func (s *Service) Create(input AuthorizationCodeInput) (*AuthorizationCode, error) {
	var a AuthorizationCode

	//TODO validation
	if input.Code != "" {
		a.Code = input.Code
	}

	if input.UserID != 0 {
		a.UserID = input.UserID
	}

	if input.AppID != 0 {
		a.AppID = input.AppID
	}

	if input.RedirectURL != "" {
		a.RedirectURL = input.RedirectURL
	}

	if input.Expiry != 0 {
		a.RedirectURL = input.RedirectURL
	}

	timeNow := time.Now().Unix()
	a.CreatedAt = timeNow
	a.UpdatedAt = timeNow

	createdAuthorizationCode, err := s.Store.Create(&a)
	if err != nil {
		return nil, err
	}

	return createdAuthorizationCode, nil
}

func (s *Service) GetByID(id string) (*AuthorizationCode, error) {

	return s.Store.GetByID(id)

}

func (s *Service) GetAll() ([]*AuthorizationCode, error) {

	return s.Store.GetAll()

}

func (s *Service) Update(id string, input *AuthorizationCodeInput) (*AuthorizationCode, error) {
	a, err := s.Store.GetByID(id)
	if err != nil {
		return nil, err
	}

	//TODO validation
	if input.Code != "" {
		a.Code = input.Code
	}

	if input.UserID != 0 {
		a.UserID = input.UserID
	}

	if input.AppID != 0 {
		a.AppID = input.AppID
	}

	if input.RedirectURL != "" {
		a.RedirectURL = input.RedirectURL
	}

	if input.Expiry != 0 {
		a.RedirectURL = input.RedirectURL
	}

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
