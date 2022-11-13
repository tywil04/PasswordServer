package main

import (
	"net/http"
	"passwordserver/src/lib/database"
	"passwordserver/src/routes"

	"github.com/joho/godotenv"
)

func handleRequests() {
	http.HandleFunc("/api/v1/auth/signin", routes.SigninHandler)
	http.HandleFunc("/api/v1/auth/signup", routes.SignupHandler)

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
