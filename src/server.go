package main

import (
	"crypto/sha512"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/pbkdf2"
)

func deriveKey(email string, masterPassword string, iterations int) []byte {
	emailBytes := []byte(email)
	masterPasswordBytes := []byte(masterPassword)
	keyLength := 512 / 8

	dk := pbkdf2.Key(emailBytes, masterPasswordBytes, iterations, keyLength, sha512.New)
	return dk
}

type SigninParams struct {
	MasterHash           string
	ProtectedDatabaseKey string
}

type SigninResponse struct {
	Errors []string
}

func homePage(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		// status := http.StatusOK
		// response := make(map[string]string)

		signinResponse := SigninResponse{}
		signinParams := SigninParams{}
		status := http.StatusOK

		decoderError := json.NewDecoder(request.Body).Decode(&signinParams)
		if decoderError != nil {
			status = http.StatusBadRequest
			signinResponse.Errors = append(signinResponse.Errors, "Unable to decode JSON body.")
		}

		if signinParams.MasterHash == "" {
			status = http.StatusBadRequest
			signinResponse.Errors = append(signinResponse.Errors, "Required parameter 'MasterHash' not provided.")
		}

		if signinParams.ProtectedDatabaseKey == "" {
			status = http.StatusBadRequest
			signinResponse.Errors = append(signinResponse.Errors, "Required parameter 'ProtectedDatabaseKey' not provided.")
		}

		if status == http.StatusOK {
			// signinResponse.Message = "This is a test message"
		}

		responseWriter.WriteHeader(status)
		responseWriter.Header().Set("Content-Type", "application/json")
		json.NewEncoder(responseWriter).Encode(&signinResponse)
	}
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8000", nil)
}

func main() {
	// fmt.Print(deriveKey("email@email.email", "masterpassword", 150000))

	handleRequests()
}
