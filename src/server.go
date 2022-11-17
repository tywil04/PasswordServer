package main

import (
	"net/http"
	customErrors "passwordserver/src/lib/cerrors"
	"passwordserver/src/lib/database"
	"passwordserver/src/routes/signin"
	"passwordserver/src/routes/signout"
	"passwordserver/src/routes/signup"
	"passwordserver/src/routes/temp"

	"github.com/joho/godotenv"
)

func handleRequests() {
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
