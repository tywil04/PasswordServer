package frontend

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed html/*
var htmlFS embed.FS
var htmlDir, _ = fs.Sub(htmlFS, "html")

func Route(handler func(response http.ResponseWriter, request *http.Request, htmlDir fs.FS)) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		handler(response, request, htmlDir)
	}
}
