package auth

import (
	// "auth/conf"
	// "auth/log"

	"net/http"
	// "time"
	// "fmt"
	"encoding/base64"
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"

	// "gorm.io/gorm"
	"github.com/go-playground/validator/v10"
)

var (
	userAuthenticated = false
	userConsent       = false
)

type AuthorizeInput struct {
	ClientID        		string	`json:"client_id"`
	UserName				string	`json:"user_name"`
	Password				string	`json:"password"`
	RedirectURL      		string  `json:"redirect_url"`  // Foreign key referencing the User table
	ResponseType       		string  `json:"response_type"`   // Foreign key referencing the App table
	Scope 					string	`json:"scope"`
	State		 			string  `json:"state"`
	CodeChallenge    		string  `json:"code_challenge"`
	CodeChallengeMethod		string  `json:"code_challenge_method"`
}

func (s *Service) CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func (s *Service) generateAuthorizationCode(length int) (code string, err error) {
	// Create a byte slice to hold random data
	randomBytes := make([]byte, length)

	// Read random data into the byte slice
	_, err = rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	// Encode the random data as a base64 string
	authorizationCode := base64.URLEncoding.EncodeToString(randomBytes)

	return authorizationCode, nil

}

func (s *Service) Authorize(input AuthorizeInput, w http.ResponseWriter, r *http.Request) (map[string]string, error) {
	var a AuthorizeInput

	userAuthenticated := false
	userConsent := false
	credentialsSent := false

	// Validate input
	err := s.Validator.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			s.Log.Error().Err(err).Msg("Validation failed")
		}
		return nil, ErrValidationFailed
	}

	if input.ClientID != "" {
		a.ClientID = input.ClientID
	}
	if input.UserName != "" {
		a.UserName = input.UserName
	}
	if input.Password != "" {
		a.Password = input.Password
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

	if a.UserName != "" && a.Password != "" && a.ClientID != ""{
		credentialsSent = true
	}
	// verify client_id/app exists
	foundApp, err := s.AuthStore.GetAppByClientID(input.ClientID)
	if err != nil {
		s.Log.Error().Err(err).Msg("Verifying that client_id exists")
		return nil, ErrLogin
	}
	s.Log.Debug().Msgf("Found client_id: ", foundApp.ClientID)


	// TODO: verify that the user belongs to app with ClientID !!!


	//TODO consent stuff !!!

	// Check if user is authenticated and has given consent
	if !userAuthenticated || !userConsent {
		if !userAuthenticated && credentialsSent {


			// Verify user_name exists
			foundUser, err := s.AuthStore.GetUserByUserName(input.UserName)
			if err != nil {
				s.Log.Error().Err(err).Msg("Verifying that user exists")
				return nil, ErrLogin
			}
			s.Log.Debug().Msgf("Found user: ", foundUser.UserName)
			// Verify hash
			match := s.CheckPasswordHash(input.Password, foundUser.Hash)
			s.Log.Debug().Msgf("Match: ", match)


			if match == false{
				s.Log.Error().Err(ErrLogin).Msg("Hashes are not equal")
				return nil, ErrLogin
			}
			s.Log.Debug().Msg("hash is equal")

			// generate authorization code
			code, err := s.generateAuthorizationCode(32)
			if err != nil {
				return nil, err
			}
			response :=map[string]string{"code": code}

			//add to db 
			authorizationCodeInput := AuthorizationCodeInput{
				Code: code,
				UserID:	foundUser.ID,
				AppID:	foundApp.ID,
			}
			createdAuthorizationCode, err := s.CreateAuthorizationCode(authorizationCodeInput)
			if err != nil {
				return nil, err
			}
			s.Log.Debug().Msgf("Created AuthorizationCode", createdAuthorizationCode)



			return response, nil
			
		}
		if !userAuthenticated && !credentialsSent {
			http.Redirect(w, r, "/v1/login?client_id=" + foundApp.ClientID, http.StatusFound)
		} else {
			http.Redirect(w, r, "/v1/login/consent", http.StatusFound)
		}
	}


	
	return nil, nil

	//TODO remove that !!!
	//response := "Hello from authorize service via get.ClientID:" + string(a.ClientID) + " RedirectURL: " + string(a.RedirectURL) + " ResponseType: " + string(a.ResponseType) + " Scope: " + string(a.Scope) + " State: " + string(a.State) + " CodeChallenge: " + string(a.CodeChallenge + " CodeChallengeMethod: " + string(a.CodeChallengeMethod))

	// return response, nil


	// Check if user is authenticated and has given consent
	// if !userAuthenticated || !userConsent {
	// 	if !userAuthenticated {
	// 		http.Redirect(w, r, "/v1/login/oauth", http.StatusFound)
	// 	} else {
	// 		http.Redirect(w, r, "/v1/login/oauth/consent", http.StatusFound)
	// 	}
	// 	return response, nil
	// }

	// Generate authorization code
	//authorizationCode := "some-authorization-code"

	// Construct redirect URL with authorization code
	//string redirectURL
	//redirectURL := fmt.Sprintf("%s?code=%s&state=%s", RedirectURL, AuthorizationCode, State)
	//var test = "http://localhost:3000/auth/login/callback?code=" + authorizationCode
	//var test = "http://localhost:3000/auth/callback"

	// Perform the redirect
	//http.Redirect(w, r, test, http.StatusFound)


}
