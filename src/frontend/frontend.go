package frontend

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
	"text/template"

	"passwordserver/src/backend"

	psPublic "passwordserver/src/lib/public"
)

//go:embed html/*
var htmlFS embed.FS
var htmlDir, _ = fs.Sub(htmlFS, "html")

var Routes map[string]string = map[string]string{}

func Route(key string, path string, handler func(response http.ResponseWriter, request *http.Request)) {
	Routes[key] = path

	http.HandleFunc(path, func(response http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet {
			handler(response, request)
		}
	})
}

type TemplateData struct {
	Data           any
	FrontendRoutes map[string]string
	BackendRoutes  map[string]string
}

func Template(patterns ...string) *template.Template {
	parts := strings.Split(patterns[0], "/")
	templateName := parts[len(parts)-1]

	return template.Must(template.New(templateName).Funcs(psPublic.Integrity).ParseFS(htmlDir, patterns...))
}

func ExecuteTemplate(response http.ResponseWriter, madeTemplate *template.Template, data any) {
	madeTemplate.Execute(response, TemplateData{Data: data, FrontendRoutes: Routes, BackendRoutes: backend.Routes})
}
