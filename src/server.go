package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"passwordserver/src/backend"
	"passwordserver/src/backend/routes"

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
	http.HandleFunc("/api/v1/auth/signin", backend.RouteHandler(http.MethodPost, routes.SigninPost))
	http.HandleFunc("/api/v1/auth/signup", backend.RouteHandler(http.MethodPost, routes.SignupPost))
	http.HandleFunc("/api/v1/auth/signout", backend.RouteHandler(http.MethodDelete, routes.SignoutDelete))
	http.HandleFunc("/temp", backend.RouteHandler(http.MethodGet, routes.TempGet))

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
