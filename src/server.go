package main

import (
	"net/http"
	"passwordserver/src/lib/database"
	"passwordserver/src/routes/signin"
	"passwordserver/src/routes/signup"

	"github.com/joho/godotenv"
)

func handleRequests() {
	http.HandleFunc("/api/v1/auth/signin", signin.SigninHandler)
	http.HandleFunc("/api/v1/auth/signup", signup.SignupHandler)

	http.ListenAndServe(":8000", nil)
}

func main() {
	dotenvError := godotenv.Load()
	if dotenvError != nil {
		panic("Error loading .env")
	}

	go database.DatabaseConnect()

	handleRequests()
}
