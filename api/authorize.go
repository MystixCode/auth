package api

import (
	// "encoding/json"
	// ut "github.com/go-playground/universal-translator"
	"auth/services/authorize"
	"auth/log"

	"encoding/json"
	"net/http"
	"fmt"
	//"github.com/gorilla/mux"
	// "github.com/go-playground/validator/v10"
	// "github.com/gorilla/mux"
)

type AuthorizeEndpoint struct {
	// logger     log.Logger
	// translator *ut.UniversalTranslator
	// validate   *validator.Validate
	service *authorize.Service
	log     *log.Logger
}

func NewAuthorizeEndpoint(log *log.Logger, service *authorize.Service) *AuthorizeEndpoint {
	return &AuthorizeEndpoint{
		service: service,
		log:     log,
	}
}

// func NewUserEndpoint(logger log.Logger, translator *ut.UniversalTranslator, validate *validator.Validate, service *user.Service) *UserEndpoint {
// 	return &UserEndpoint{
// 		logger:     logger.WithPrefix("api.user"),
// 		translator: translator,
// 		validate:   validate,
// 		service:    service,
// 	}
// }


func (e *AuthorizeEndpoint) ViaPost(w http.ResponseWriter, r *http.Request) {
	var input authorize.AuthorizeInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respond(w, e.log, http.StatusBadRequest, "invalid body", nil)
		return
	}


	response, err := e.service.ViaPost(input, w, r)
	if err != nil {
		respond(w, e.log, http.StatusBadRequest, "invalid body", nil)
	}

	respond(w, e.log, http.StatusOK, "ViaPost", response)
}

func (e *AuthorizeEndpoint) ViaGet(w http.ResponseWriter, r *http.Request) {
	var input authorize.AuthorizeInput

	
	input.ClientID = r.URL.Query().Get("client_id")




	response, err := e.service.ViaGet(input, w, r)
	if err != nil {
		respond(w, e.log, http.StatusBadRequest, "invalid body", nil)
	}


	respond(w, e.log, http.StatusOK, "ViaGet", response)
}

func (e *AuthorizeEndpoint) LoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//TODO das ganze zum service durchschlaufe services/authorize/service.go
	loginForm := `
	<form method="post" action="/v1/login">
		<input type="text" name="username" placeholder="Username"></br>
		<input type="password" name="password" placeholder="Password"></br>
		<button type="submit">Login</button>
	</form>
`
	fmt.Fprintln(w, loginForm)

	//respond(w, e.log, http.StatusOK, fmt.Fprintln(w, loginForm), nil)
}

func (e *AuthorizeEndpoint) Login(w http.ResponseWriter, r *http.Request) {

	respond(w, e.log, http.StatusOK, "Todo: Login", nil)
}

func (e *AuthorizeEndpoint) Consent(w http.ResponseWriter, r *http.Request) {

	respond(w, e.log, http.StatusOK, "Todo: Consent", nil)
}