package api

import (
	// "encoding/json"
	// ut "github.com/go-playground/universal-translator"
	"auth/services/auth"
	"auth/log"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	// "github.com/go-playground/validator/v10"
	// "github.com/gorilla/mux"
)

type AuthorizationCodeEndpoint struct {
	// logger     log.Logger
	// translator *ut.UniversalTranslator
	// validate   *validator.Validate
	service *auth.Service
	log     *log.Logger
}

func NewAuthorizationCodeEndpoint(log *log.Logger, service *auth.Service) *AuthorizationCodeEndpoint {
	return &AuthorizationCodeEndpoint{
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

func (e *AuthorizationCodeEndpoint) Create(w http.ResponseWriter, r *http.Request) {
	var input auth.AuthorizationCodeInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respond(w, e.log, http.StatusBadRequest, "invalid body", nil)
		return
	}
	//fmt.Println("input", input)
	// err = e.validate.Struct(input)
	// if err != nil {
	// 	errs := getValidationError(err.(validator.ValidationErrors), trans)
	// 	respond(w, e.logger, http.StatusBadRequest, "validation failed", errs)
	// 	return
	// }

	createdAuthorizationCode, err := e.service.CreateAuthorizationCode(input)
	if err != nil {
		respond(w, e.log, http.StatusBadRequest, "invalid body", nil)
	}

	respond(w, e.log, http.StatusCreated, "AuthorizationCode created successfully", createdAuthorizationCode)
}

func (e *AuthorizationCodeEndpoint) GetAll(w http.ResponseWriter, _ *http.Request) {
	authorizationCodes, err := e.service.GetAllAuthorizationCodes()
	if err != nil {
		switch err {
		case auth.ErrNotFound:
			respond(w, e.log, http.StatusNotFound, err.Error(), nil)
		default:
			respond(w, e.log, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	respond(w, e.log, http.StatusOK, "load all AuthorizationCodes successfully", authorizationCodes)
}

func (e *AuthorizationCodeEndpoint) GetById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	foundAuthorizationCode, err := e.service.GetAuthorizationCodeByID(id)
	if err != nil {
		switch err {
		case auth.ErrIdParseFailed:
			respond(w, e.log, http.StatusBadRequest, err.Error(), nil)
		case auth.ErrNotFound:
			respond(w, e.log, http.StatusNotFound, err.Error(), nil)
		default:
			respond(w, e.log, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	respond(w, e.log, http.StatusOK, "Successfully found AuthorizationCode", foundAuthorizationCode)
}

func (e *AuthorizationCodeEndpoint) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := e.service.DeleteAuthorizationCode(id)
	if err != nil {
		switch err {
		case auth.ErrIdParseFailed:
			respond(w, e.log, http.StatusBadRequest, err.Error(), nil)
		case auth.ErrNotFound:
			respond(w, e.log, http.StatusNotFound, err.Error(), nil)
		default:
			respond(w, e.log, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	respond(w, e.log, http.StatusOK, "successfully deleted", nil)
}
