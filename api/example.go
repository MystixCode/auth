package api

import (
	// "encoding/json"
	// ut "github.com/go-playground/universal-translator"
	"auth/log"

	"encoding/json"
	"net/http"

	"auth/services/example"

	"github.com/gorilla/mux"
	// "github.com/go-playground/validator/v10"
	// "github.com/gorilla/mux"
)

type ExampleEndpoint struct {
	log *log.Logger
	// translator *ut.UniversalTranslator
	// validate   *validator.Validate
	service *example.Service
}

func NewExampleEndpoint(log *log.Logger, service *example.Service) *ExampleEndpoint {
	return &ExampleEndpoint{
		service: service,
		log:     log,
	}
}

// func NewExampleEndpoint(logger log.Logger, translator *ut.UniversalTranslator, validate *validator.Validate, service *example.Service) *ExampleEndpoint {
// 	return &ExampleEndpoint{
// 		logger:     logger.WithPrefix("api.example"),
// 		translator: translator,
// 		validate:   validate,
// 		service:    service,
// 	}
// }

func (e *ExampleEndpoint) Create(w http.ResponseWriter, r *http.Request) {
	var input example.ExampleInput

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

	createdExample, err := e.service.Create(input)
	if err != nil {
		respond(w, e.log, http.StatusBadRequest, "invalid body", nil)
	}

	respond(w, e.log, http.StatusCreated, "example created successfully", createdExample)
}

func (e *ExampleEndpoint) GetAll(w http.ResponseWriter, _ *http.Request) {
	examples, err := e.service.GetAll()
	if err != nil {
		switch err {
		case example.ErrNotFound:
			respond(w, e.log, http.StatusNotFound, err.Error(), nil)
		default:
			respond(w, e.log, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	respond(w, e.log, http.StatusOK, "load all examples successfully", examples)
}

func (e *ExampleEndpoint) GetById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	foundExample, err := e.service.GetByID(id)
	if err != nil {
		switch err {
		case example.ErrIdParseFailed:
			respond(w, e.log, http.StatusBadRequest, err.Error(), nil)
		case example.ErrNotFound:
			respond(w, e.log, http.StatusNotFound, err.Error(), nil)
		default:
			respond(w, e.log, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	respond(w, e.log, http.StatusOK, "successfully found example", foundExample)
}

func (e *ExampleEndpoint) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var input *example.ExampleInput
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

	createdExample, err := e.service.Update(id, input)
	if err != nil {
		switch err {
		case example.ErrIdParseFailed:
			respond(w, e.log, http.StatusBadRequest, err.Error(), nil)
		case example.ErrNotFound:
			respond(w, e.log, http.StatusNotFound, err.Error(), nil)
		default:
			respond(w, e.log, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	respond(w, e.log, http.StatusCreated, "example updated successfully", createdExample)
}

func (e *ExampleEndpoint) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := e.service.Delete(id)
	if err != nil {
		switch err {
		case example.ErrIdParseFailed:
			respond(w, e.log, http.StatusBadRequest, err.Error(), nil)
		case example.ErrNotFound:
			respond(w, e.log, http.StatusNotFound, err.Error(), nil)
		default:
			respond(w, e.log, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	respond(w, e.log, http.StatusOK, "successfully deleted", nil)
}
