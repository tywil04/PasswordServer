package pages

import (
	"net/http"
	"passwordserver/src/frontend"
)

type SigninData struct {
	Error string
}

func Signin(response http.ResponseWriter, request *http.Request) {
	template := frontend.Template("auth/signin.html", "base.html")
	frontend.ExecuteTemplate(response, template, SigninData{})
}
