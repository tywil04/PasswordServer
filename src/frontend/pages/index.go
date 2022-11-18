package pages

import (
	"io/fs"
	"net/http"
	"text/template"
)

type IndexParameters struct {
	Test string
}

func IndexHandler(response http.ResponseWriter, request *http.Request, htmlDir fs.FS) {
	template := template.Must(template.ParseFS(htmlDir, "index.html"))

	indexParameters := IndexParameters{Test: "Testing Context!"}

	template.Execute(response, indexParameters)
}
