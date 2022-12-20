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
	psFS "passwordserver/src/lib/fs"
	psPublic "passwordserver/src/lib/public"

	"github.com/joho/godotenv"
)

//go:embed public/*
var publicFS embed.FS
var publicDir, _ = fs.Sub(publicFS, "public")

func handleRequests() {
	// Static Path
	http.Handle("/", http.StripPrefix(strings.TrimRight("/public", "/"), http.FileServer(psFS.FileSystem{Fs: http.FS(publicDir)})))

	// API Routes
	backend.Route("signin", "/api/v1/auth/signin", backend.MethodMap{Post: routes.SigninPost})
	backend.Route("signup", "/api/v1/auth/signup", backend.MethodMap{Post: routes.SignupPost})
	backend.Route("signout", "/api/v1/auth/signout", backend.MethodMap{Delete: routes.SignoutDelete})
	backend.Route("protectedDatabaseKey", "/api/v1/user/protectedDatabaseKey", backend.MethodMap{Get: routes.ProtectedDatabaseKeyGet})

	// Pages
	frontend.Route("signin", "/auth/signin", pages.Signin)
	frontend.Route("signup", "/auth/signup", pages.Signup)
	frontend.Route("home", "/home", pages.Home)

	http.ListenAndServe(":8000", nil)
}

func main() {
	dotenvError := godotenv.Load()
	if dotenvError != nil {
		panic(psErrors.ErrorLoadingEnv)
	}

	database := psDatabase.DatabaseConnect()
	if database == nil {
		panic(psErrors.ErrorLoadingDatabase)
	}

	psPublic.GenerateIntegrityMap(publicDir)

	handleRequests()
}
