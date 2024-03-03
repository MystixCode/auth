package api

import (
	"auth/log"

	"encoding/json"
	"net/http"

	// ut "github.com/go-playground/universal-translator"
	// "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Api struct {
	Root   				*RootEndpoint
	Health 				*HealthEndpoint
	User    			*UserEndpoint
	App    				*AppEndpoint
	AuthorizationCode	*AuthorizationCodeEndpoint
	Login				*LoginEndpoint
	Example				*ExampleEndpoint
}

type Body struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (a *Api) New(router *mux.Router) {
	v1 := router.PathPrefix("/v1").Subrouter().StrictSlash(false)

	// Root
	v1.HandleFunc("", a.Root.GetRoot).Methods(http.MethodGet)
	// v1.HandleFunc("", a.Root.GetRoot).Methods(http.MethodGet)

	// Health
	v1.HandleFunc("/health", a.Health.GetHealth).Methods(http.MethodGet)

	// Example
	v1.HandleFunc("/examples", a.Example.Create).Methods(http.MethodPost)
	v1.HandleFunc("/examples", a.Example.GetAll).Methods(http.MethodGet)
	v1.HandleFunc("/examples/{id}", a.Example.GetById).Methods(http.MethodGet)
	v1.HandleFunc("/examples/{id}", a.Example.Update).Methods(http.MethodPut)
	v1.HandleFunc("/examples/{id}", a.Example.Delete).Methods(http.MethodDelete)

	// User
	v1.HandleFunc("/users", a.User.Create).Methods(http.MethodPost)
	v1.HandleFunc("/users", a.User.GetAll).Methods(http.MethodGet)
	v1.HandleFunc("/users/{id}", a.User.GetById).Methods(http.MethodGet)
	v1.HandleFunc("/users/{id}", a.User.Update).Methods(http.MethodPut)
	v1.HandleFunc("/users/{id}", a.User.Delete).Methods(http.MethodDelete)

	// App
	v1.HandleFunc("/apps", a.App.Create).Methods(http.MethodPost)
	v1.HandleFunc("/apps", a.App.GetAll).Methods(http.MethodGet)
	v1.HandleFunc("/apps/{id}", a.App.GetById).Methods(http.MethodGet)
	v1.HandleFunc("/apps/{id}", a.App.Update).Methods(http.MethodPut)
	v1.HandleFunc("/apps/{id}", a.App.Delete).Methods(http.MethodDelete)

	// AuthorizationCode
	v1.HandleFunc("/authorizationcodes", a.AuthorizationCode.Create).Methods(http.MethodPost)
	v1.HandleFunc("/authorizationcodes", a.AuthorizationCode.GetAll).Methods(http.MethodGet)
	v1.HandleFunc("/authorizationcodes/{id}", a.AuthorizationCode.GetById).Methods(http.MethodGet)
	v1.HandleFunc("/authorizationcodes/{id}", a.AuthorizationCode.Delete).Methods(http.MethodDelete)

	// Oauth Flows endpoints
	v1.HandleFunc("/login/oauth/authorize", a.Login.Authorize).Methods(http.MethodGet)
	v1.HandleFunc("/login/oauth/authorize", a.Login.Authorize).Methods(http.MethodPost)
	v1.HandleFunc("/login", a.Login.LoginPage).Methods(http.MethodGet)
	v1.HandleFunc("/login/token", a.Login.Token).Methods(http.MethodPost)
	// v1.HandleFunc("/login/oauth/consent", a.Login.Consent).Methods(http.MethodGet)

}

func respond(w http.ResponseWriter, log *log.Logger, status int, message string, data interface{}) {
	body := Body{}
	body.Status = status
	w.WriteHeader(status)
	if message != "" {
		body.Message = message
	}

	if data != nil {
		body.Data = data
	}

	bodyByte, err := json.Marshal(body)
	if err != nil {
		log.Error().Err(err).Msg("fail to parse response body")
		respond(w, log, http.StatusInternalServerError, "internal error", nil)
		return
	}

	headerByte, err := json.Marshal(w.Header())
	if err != nil {
		log.Error().Err(err).Msg("fail to parse response header")
		respond(w, log, http.StatusInternalServerError, "internal error", nil)
		return
	}

	log.Debug().RawJSON("header", headerByte).RawJSON("body", bodyByte).Msg("Response")

	_, err = w.Write(bodyByte)
	if err != nil {
		log.Error().Err(err).Msg("fail to write response body")
		respond(w, log, http.StatusInternalServerError, "internal error", nil)
		return
	}
}

// func getTranslator(r *http.Request, translator *ut.UniversalTranslator) ut.Translator {
// 	lang := r.Header.Get("Accept-Language")
// 	trans, _ := translator.GetTranslator(lang)
// 	return trans
// }

// func getValidationError(err validator.ValidationErrors, translator ut.Translator) ValidationErrors {
// 	ve := ValidationErrors{}

// 	for _, e := range err {
// 		ve = append(ve, ValidationError{e.Field(): e.Translate(translator)})
// 	}

// 	return ve
// }
