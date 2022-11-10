package main

import (
	"crypto/sha512"
	"net/http"
	"passwordserver/src/routes"

	"golang.org/x/crypto/pbkdf2"
)

func deriveKey(email string, masterPassword string, iterations int) []byte {
	emailBytes := []byte(email)
	masterPasswordBytes := []byte(masterPassword)
	keyLength := 512 / 8

	dk := pbkdf2.Key(emailBytes, masterPasswordBytes, iterations, keyLength, sha512.New)
	return dk
}

func handleRequests() {
	http.HandleFunc("/api/v1/auth/signin", routes.SigninHandler)
	http.HandleFunc("/api/v1/auth/signup", routes.SignupHandler)

	http.ListenAndServe(":8000", nil)
}

func main() {
	// fmt.Print(deriveKey("email@email.email", "masterpassword", 150000))

	handleRequests()
}
