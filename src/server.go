package main

import (
	"embed"
	"io/fs"
	"net/http"
	"passwordserver/src/backend/routes/signin"
	"passwordserver/src/backend/routes/signout"
	"passwordserver/src/backend/routes/signup"
	"passwordserver/src/backend/routes/temp"
	customErrors "passwordserver/src/lib/cerrors"
	"passwordserver/src/lib/database"
	customFS "passwordserver/src/lib/fs"
	"strings"

	"github.com/joho/godotenv"
)

//go:embed public/*
var publicFS embed.FS
var publicDir, _ = fs.Sub(publicFS, "public")

func handleRequests() {
	http.Handle("/", http.StripPrefix(strings.TrimRight("/public", "/"), http.FileServer(customFS.FileSystem{Fs: http.FS(publicDir)})))
	http.HandleFunc("/api/v1/auth/signin", signin.SigninHandler)
	http.HandleFunc("/api/v1/auth/signup", signup.SignupHandler)
	http.HandleFunc("/api/v1/auth/signout", signout.SignoutHandler)
	http.HandleFunc("/temp", temp.TempHandler)

	http.ListenAndServe(":8000", nil)
}

func main() {
	dotenvError := godotenv.Load()
	if dotenvError != nil {
		panic(customErrors.ErrorLoadingEnv)
	}

	go database.DatabaseConnect()

	handleRequests()
}
