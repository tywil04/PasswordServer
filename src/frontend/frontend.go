package frontend

import (
	"embed"
	"io/fs"
	"net/http"
	"text/template"
)

//go:embed html/*
var htmlFS embed.FS
var htmlDir, _ = fs.Sub(htmlFS, "html")

func RouteHandler(htmlPath string, handler func(response http.ResponseWriter, request *http.Request, template *template.Template)) func(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFS(htmlDir, htmlPath, "layout.html"))
	return func(response http.ResponseWriter, request *http.Request) {
		handler(response, request, template)
	}
}
