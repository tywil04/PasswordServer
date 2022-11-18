package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"passwordserver/src/backend"
	"passwordserver/src/backend/routes"
	"passwordserver/src/frontend"
	"passwordserver/src/frontend/pages"

	psDatabase "passwordserver/src/lib/database"
	psErrors "passwordserver/src/lib/errors"
	psCustomFS "passwordserver/src/lib/fs"

	"github.com/joho/godotenv"
)

//go:embed public/*
var publicFS embed.FS
var publicDir, _ = fs.Sub(publicFS, "public")

func handleRequests() {
	// Static Path
	http.Handle("/", http.StripPrefix(strings.TrimRight("/public", "/"), http.FileServer(psCustomFS.FileSystem{Fs: http.FS(publicDir)})))

	// API Routes
	http.HandleFunc("/api/v1/auth/signin", backend.Route(backend.MethodMap{Post: routes.SigninPost}))
	http.HandleFunc("/api/v1/auth/signup", backend.Route(backend.MethodMap{Post: routes.SignupPost}))
	http.HandleFunc("/api/v1/auth/signout", backend.Route(backend.MethodMap{Delete: routes.SignoutDelete}))
	http.HandleFunc("/temp", backend.Route(backend.MethodMap{Get: routes.TempGet}))

	// Pages
	http.HandleFunc("/testing", frontend.RouteHandler("index.html", pages.IndexHandler))

	http.ListenAndServe(":8000", nil)
}

func main() {
	dotenvError := godotenv.Load()
	if dotenvError != nil {
		panic(psErrors.ErrorLoadingEnv)
	}

	go psDatabase.DatabaseConnect()

	handleRequests()
}
