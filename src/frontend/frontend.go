package frontend

import (
	"bytes"
	"crypto/sha512"
	"embed"
	"encoding/base64"
	"io/fs"
	"net/http"
	"strings"
	"text/template"

	"passwordserver/src/backend"
)

//go:embed html/*
var htmlFS embed.FS
var htmlDir, _ = fs.Sub(htmlFS, "html")

//go:embed js/*
var jsFS embed.FS
var jsDir, _ = fs.Sub(jsFS, "js")

//go:embed css/*
var cssFS embed.FS
var cssDir, _ = fs.Sub(cssFS, "css")

var Routes map[string]string = map[string]string{}
var RoutesJS map[string]string = map[string]string{}
var RoutesCSS map[string]string = map[string]string{}

func Route(key string, path string, handler func(response http.ResponseWriter, request *http.Request)) {
	Routes[key] = path

	http.HandleFunc(path, func(response http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet {
			handler(response, request)
		}
	})
}

func RouteCSS(key string, path string, filePath string) {
	RoutesCSS[key] = path
	file, _ := fs.ReadFile(cssDir, filePath)

	http.HandleFunc(path, func(response http.ResponseWriter, request *http.Request) {
		response.Header().Add("Content-type", "text/css")
		response.Write(file)
	})
}

var CachedIntegrity map[string]string = map[string]string{}
var Integrity template.FuncMap = template.FuncMap{}

func RouteJS(key string, path string, filePath string) {
	RoutesJS[key] = path
	file, _ := fs.ReadFile(jsDir, filePath)

	http.HandleFunc(path, func(response http.ResponseWriter, request *http.Request) {
		response.Header().Add("Content-type", "text/javascript")
		template.Must(template.New("").Parse(bytes.NewBuffer(file).String())).Execute(response, TemplateData{FrontendRoutes: Routes, FrontendRoutesJS: RoutesJS, FrontendRoutesCSS: RoutesCSS, BackendRoutes: backend.Routes})
	})

	key += "PublicIntegrity"

	Integrity[key] = func() string {
		if CachedIntegrity[key] == "" {
			hashAlgo := "sha384"
			hash := sha512.New384()

			contents, _ := fs.ReadFile(jsDir, filePath)
			parsedContents := bytes.NewBuffer([]byte{})
			template.Must(template.New("").Parse(bytes.NewBuffer(contents).String())).Execute(parsedContents, TemplateData{FrontendRoutes: Routes, FrontendRoutesJS: RoutesJS, FrontendRoutesCSS: RoutesCSS, BackendRoutes: backend.Routes})
			hash.Write(parsedContents.Bytes())

			sum := hash.Sum(nil)
			resultBuffer := bytes.NewBuffer([]byte{})
			base64.NewEncoder(base64.StdEncoding, resultBuffer).Write(sum)
			value := hashAlgo + "-" + resultBuffer.String()

			CachedIntegrity[key] = value

			return value
		} else {
			return CachedIntegrity[key]
		}
	}
}

type TemplateData struct {
	Data              any
	FrontendRoutes    map[string]string
	FrontendRoutesJS  map[string]string
	FrontendRoutesCSS map[string]string
	BackendRoutes     map[string]string
}

func Template(patterns ...string) *template.Template {
	parts := strings.Split(patterns[0], "/")
	templateName := parts[len(parts)-1]

	return template.Must(template.New(templateName).Funcs(Integrity).ParseFS(htmlDir, patterns...))
}

func ExecuteTemplate(response http.ResponseWriter, madeTemplate *template.Template, data any) {
	madeTemplate.Execute(response, TemplateData{Data: data, FrontendRoutes: Routes, FrontendRoutesJS: RoutesJS, FrontendRoutesCSS: RoutesCSS, BackendRoutes: backend.Routes})
}
