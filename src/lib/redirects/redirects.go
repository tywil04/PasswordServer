package redirects

import (
	"net/http"
	"passwordserver/src/frontend"
)

func RedirectSignin(response http.ResponseWriter, status int) {
	response.Header().Set("Location", frontend.Routes["signin"])
	response.WriteHeader(status)
}

func RedirectSignup(response http.ResponseWriter, status int) {
	response.Header().Set("Location", frontend.Routes["signup"])
	response.WriteHeader(status)
}

func RedirectHome(response http.ResponseWriter, status int) {
	response.Header().Set("Location", frontend.Routes["home"])
	response.WriteHeader(status)
}
