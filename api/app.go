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

type AppEndpoint struct {
	// logger     log.Logger
	// translator *ut.UniversalTranslator
	// validate   *validator.Validate
	service *auth.Service
	log     *log.Logger
}

func NewAppEndpoint(log *log.Logger, service *auth.Service) *AppEndpoint {
	return &AppEndpoint{
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

func (e *AppEndpoint) Create(w http.ResponseWriter, r *http.Request) {
	var input auth.AppInput

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

	createdApp, err := e.service.CreateApp(input)
	if err != nil {
		respond(w, e.log, http.StatusBadRequest, "invalid body", nil)
	}

	respond(w, e.log, http.StatusCreated, "App created successfully", createdApp)
}

func (e *AppEndpoint) GetAll(w http.ResponseWriter, _ *http.Request) {
	apps, err := e.service.GetAllApps()
	if err != nil {
		switch err {
		case auth.ErrNotFound:
			respond(w, e.log, http.StatusNotFound, err.Error(), nil)
		default:
			respond(w, e.log, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	respond(w, e.log, http.StatusOK, "load all Apps successfully", apps)
}

func (e *AppEndpoint) GetById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	foundApp, err := e.service.GetAppByID(id)
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

	respond(w, e.log, http.StatusOK, "successfully found App", foundApp)
}

func (e *AppEndpoint) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var input *auth.AppInput
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

	createdApp, err := e.service.UpdateApp(id, input)
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

	respond(w, e.log, http.StatusCreated, "App updated successfully", createdApp)
}

func (e *AppEndpoint) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := e.service.DeleteApp(id)
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
