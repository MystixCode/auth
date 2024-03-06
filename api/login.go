package api

import (
	// "encoding/json"
	// ut "github.com/go-playground/universal-translator"
	"auth/log"
	"auth/services/auth"

	"encoding/json"
	"fmt"
	"net/http"
	// "github.com/gorilla/mux"
	// "github.com/go-playground/validator/v10"
	// "github.com/gorilla/mux"
)

type LoginEndpoint struct {
	// logger     log.Logger
	// translator *ut.UniversalTranslator
	// validate   *validator.Validate
	service *auth.Service
	log     *log.Logger
}

func NewLoginEndpoint(log *log.Logger, service *auth.Service) *LoginEndpoint {
	return &LoginEndpoint{
		service: service,
		log:     log,
	}
}

func (e *LoginEndpoint) Authorize(w http.ResponseWriter, r *http.Request) {
	var input auth.AuthorizeInput
	e.log.Debug().Msg("Authorize__________________________________________________")

	// parse input
	if r.Method == http.MethodGet {
		input.ClientID = r.URL.Query().Get("client_id")
		//input.username = r.URL.Query().Get("client_id")
	} else {
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			respond(w, e.log, http.StatusBadRequest, "invalid body", nil)
			return
		}
	}

	response, err := e.service.Authorize(input, w, r)
	if err != nil {
		respond(w, e.log, http.StatusBadRequest, "invalid body", nil)
	}

	respond(w, e.log, http.StatusOK, "ViaPost", response)
}

func (e *LoginEndpoint) LoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//TODO das ganze zum service durchschlaufe services/auth/authorize.go oder sogar in a eigene service wegisoliere !!!

	loginForm := `
	<html>
	<head>
	</head>
	<body>
        <style>
            .mainContainer {
                min-height: 100vh;
                font-family: Helvetica, Arial, sans-serif;
                color: #eceff4;
                background-color:#121212;
                padding: 24px;
            }
            .authContainer {
                display: fley;
                flex-direction: column;
                align-items: center;
                justify-content: center;
                flex: 1;
                padding: 24px;
                background-color: #282828;
            }

            * {
                font-family: arial;
                box-sizing: border-box;
                padding: 0;
                margin: 0;
            }

            form {
                display: flex;
                flex-direction: column;
                gap: 8px;
                width: 400px;
                max-width: 100%;
                margin: 0 auto;
            }

            form input {
                width: 100%;
                font-size: 12px;
                padding: 4px;

            }

            button {
                width: 100%;
                background-color: #0CAF50; /* Green */
                border: none;
                color: white;
                font-size: 16px;
                padding: 4px;
                text-align: center;
                text-decoration: none;
                display: inline-block;
            }
            button:hover {
                background-color: #4CAF59; /* Green */
            }

            h1 {
                text-align: center;
            }
            a:link {
                text-decoration: none;
                color: #ccc;
            }

            a:visited {
                text-decoration: none;
                color: #ccc;
            }

            a:hover {
                text-decoration: none;
                color: #ccc;
            }

            a:active { text-decoration: none; }
        </style>


        <div class="mainContainer">
            <div class="authContainer">
                <form>
                    <h1>Login</h1>
                    <label>
                        <input type="text" name="user_name" id="user_name" placeholder="Username" />
                    </label>
                    <label>
                        <input type="password" name="password" id="password" placeholder="Password" />
                    </label>
                    <label>
                        <button type="submit">Login</button>
                        </br>
                        <a href="http://localhost:3000/signup">Signup</a> | <a href="http://localhost:3000/forgot">Forgot Password</a>
                    </label>
                </form>
            </div>
        </div>

        <script type="application/javascript">
            function handleSubmit(event) {
                event.preventDefault();

                const data = new FormData(event.target);

                const jsonData = Object.fromEntries(data.entries());
                console.log("---------------------------------");
                console.log({ jsonData });

                const queryString = window.location.search;
                console.log(queryString);
                const urlParams = new URLSearchParams(queryString);
                if (urlParams.has('client_id')) {
                    const clientId = urlParams.get('client_id')
                    console.log(clientId);
                    jsonData.client_id = clientId ;
                }

                console.log("---------------------------------");
                console.log({ jsonData });

                fetch('http://localhost:8080/v1/login/oauth/authorize', {
                    method: 'POST',
                    headers: {
                        'Accept': 'application/json',
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(jsonData)
                })
                    .then(response => response.json())
                    .then(response => window.location.replace('http://localhost:3000/auth/callback?code=' + response.data.code))
            }

            const form = document.querySelector('form');
            form.addEventListener('submit', handleSubmit);
        </script>

    </body>
</html>

`

	// .then(response => console.log(JSON.stringify(response)))
	fmt.Fprintln(w, loginForm)

	//respond(w, e.log, http.StatusOK, fmt.Fprintln(w, loginForm), nil)
}

func (e *LoginEndpoint) Token(w http.ResponseWriter, r *http.Request) {
	var input auth.TokenInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respond(w, e.log, http.StatusBadRequest, "invalid body", nil)
		return
	}

	code, err := e.service.Token(input)
	if err != nil {
		respond(w, e.log, http.StatusBadRequest, err.Error(), nil)
		return
	}

	respond(w, e.log, http.StatusOK, "user logged in successfully", code)
}

func (e *LoginEndpoint) Consent(w http.ResponseWriter, r *http.Request) {

	respond(w, e.log, http.StatusOK, "Todo: Consent", nil)
}
