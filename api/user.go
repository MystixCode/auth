package api

import (
	// "encoding/json"
	// ut "github.com/go-playground/universal-translator"
	"auth/services/user"
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
	service *user.Service
	log     *log.Logger
}

func NewUserEndpoint(log *log.Logger, service *user.Service) *UserEndpoint {
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
	var input user.LoginInput

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
	var input user.UserInput

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

	createdUser, err := e.service.Create(input)
	if err != nil {
		respond(w, e.log, http.StatusBadRequest, "invalid body", nil)
		return
	}

	respond(w, e.log, http.StatusCreated, "user created successfully", createdUser)
}

func (e *UserEndpoint) GetAll(w http.ResponseWriter, _ *http.Request) {
	users, err := e.service.GetAll()
	if err != nil {
		switch err {
		case user.ErrNotFound:
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

	foundUser, err := e.service.GetByID(id)
	if err != nil {
		switch err {
		case user.ErrIdParseFailed:
			respond(w, e.log, http.StatusBadRequest, err.Error(), nil)
		case user.ErrNotFound:
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
	var input *user.UserInput
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

	createdUser, err := e.service.Update(id, input)
	if err != nil {
		switch err {
		case user.ErrIdParseFailed:
			respond(w, e.log, http.StatusBadRequest, err.Error(), nil)
		case user.ErrNotFound:
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

	err := e.service.Delete(id)
	if err != nil {
		switch err {
		case user.ErrIdParseFailed:
			respond(w, e.log, http.StatusBadRequest, err.Error(), nil)
		case user.ErrNotFound:
			respond(w, e.log, http.StatusNotFound, err.Error(), nil)
		default:
			respond(w, e.log, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	respond(w, e.log, http.StatusOK, "successfully deleted", nil)
}
