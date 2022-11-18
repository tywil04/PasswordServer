package pages

import (
	"net/http"
	"text/template"
)

type IndexParameters struct {
	Test string
}

func IndexHandler(response http.ResponseWriter, request *http.Request, template *template.Template) {
	indexParameters := IndexParameters{Test: "Testing Context!"}
	template.Execute(response, indexParameters)
}
