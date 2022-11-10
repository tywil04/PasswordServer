package main

import (
	"net/http"
	"passwordserver/src/lib"
	"passwordserver/src/routes"
)

func handleRequests() {
	http.HandleFunc("/api/v1/auth/signin", routes.SigninHandler)
	http.HandleFunc("/api/v1/auth/signup", routes.SignupHandler)

	http.ListenAndServe(":8000", nil)
}

func main() {
	go lib.DatabaseConnect()
	handleRequests()
}
