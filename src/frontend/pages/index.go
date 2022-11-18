package pages

import (
	"io/fs"
	"net/http"
	"text/template"
)

type IndexData struct {
	Test string
}

func IndexHandler(response http.ResponseWriter, request *http.Request, htmlDir fs.FS) {
	template := template.Must(template.ParseFS(htmlDir, "index.html"))

	indexData := IndexData{Test: "Testing Context!"}

	template.Execute(response, indexData)
}
