package redirect

import "net/http"

func RedirectSignin(response http.ResponseWriter, status int) {
	response.Header().Set("Location", "/auth/signin")
	response.WriteHeader(status)
}

func RedirectSignup(response http.ResponseWriter, status int) {
	response.Header().Set("Location", "/auth/signup")
	response.WriteHeader(status)
}

func RedirectHome(response http.ResponseWriter, status int) {
	response.Header().Set("Location", "/home")
	response.WriteHeader(status)
}
