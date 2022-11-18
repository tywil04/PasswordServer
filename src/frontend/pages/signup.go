package pages

import (
	"io/fs"
	"net/http"
	"text/template"
)

type SignupData struct {
	Error string
}

func Signup(response http.ResponseWriter, request *http.Request, htmlDir fs.FS) {
	template := template.Must(template.ParseFS(htmlDir, "signup.html", "layout.html"))

	template.Execute(response, SignupData{})
}
