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

type UserEndpoint struct {
	// logger     log.Logger
	// translator *ut.UniversalTranslator
	// validate   *validator.Validate
	service *auth.Service
	log     *log.Logger
}

func NewUserEndpoint(log *log.Logger, service *auth.Service) *UserEndpoint {
	return &UserEndpoint{
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

func (e *UserEndpoint) Login(w http.ResponseWriter, r *http.Request) {
	var input auth.LoginInput

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

	response, err := e.service.Login(input)
	if err != nil {
		respond(w, e.log, http.StatusBadRequest, err.Error(), nil)
		return
	}

	respond(w, e.log, http.StatusOK, "user logged in successfully", response)
}

func (e *UserEndpoint) Create(w http.ResponseWriter, r *http.Request) {
	var input auth.UserInput

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

	createdUser, err := e.service.CreateUser(input)
	if err != nil {
		respond(w, e.log, http.StatusBadRequest, "invalid body", nil)
		return
	}

	respond(w, e.log, http.StatusCreated, "user created successfully", createdUser)
}

func (e *UserEndpoint) GetAll(w http.ResponseWriter, _ *http.Request) {
	users, err := e.service.GetAllUsers()
	if err != nil {
		switch err {
		case auth.ErrNotFound:
			respond(w, e.log, http.StatusNotFound, err.Error(), nil)
		default:
			respond(w, e.log, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	respond(w, e.log, http.StatusOK, "load all users successfully", users)
}

func (e *UserEndpoint) GetById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	foundUser, err := e.service.GetUserByID(id)
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

	respond(w, e.log, http.StatusOK, "successfully found user", foundUser)
}

func (e *UserEndpoint) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var input *auth.UserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respond(w, e.log, http.StatusBadRequest, "invalid body", nil)
		return
	}

	// err = e.validate.Struct(input)
	// if err != nil {
	// 	errs := getValidationError(err.(validator.ValidationErrors), trans)
	// 	respond(w, e.logger, http.StatusBadRequest, "validation failed", errs)
	// 	return
	// }

	createdUser, err := e.service.UpdateUser(id, input)
	if err != nil {
		switch err {
		case auth.ErrIdParseFailed:
			respond(w, e.log, http.StatusBadRequest, err.Error(), nil)
		case auth.ErrNotFound:
			respond(w, e.log, http.StatusNotFound, err.Error(), nil)
		// case user.ErrPasswordChangeNotAllowed:
		// 	respond(w, http.StatusBadRequest, err.Error(), nil)
		default:
			respond(w, e.log, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	respond(w, e.log, http.StatusCreated, "user updated successfully", createdUser)
}

func (e *UserEndpoint) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := e.service.DeleteUser(id)
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
