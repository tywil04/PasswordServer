package pages

import (
	"io/fs"
	"net/http"
	"text/template"
)

type SigninData struct {
	Error string
}

func Signin(response http.ResponseWriter, request *http.Request, htmlDir fs.FS) {
	template := template.Must(template.ParseFS(htmlDir, "auth/signin.html", "base.html"))

	template.Execute(response, SigninData{})
}
