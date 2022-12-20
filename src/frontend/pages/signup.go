package pages

import (
	"net/http"
	"passwordserver/src/frontend"
)

type SignupData struct {
	Error string
}

func Signup(response http.ResponseWriter, request *http.Request) {
	template := frontend.Template("signup.html", "base.html")
	frontend.ExecuteTemplate(response, template, SignupData{})
}
