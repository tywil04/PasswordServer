package signin

import (
	"encoding/json"
	"net/http"
)

type SigninParams struct {
	MasterHash           string
	ProtectedDatabaseKey string
}

type SigninResponse struct {
	Errors []string
}

func Handler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		post(response, request)
	}
}

func post(response http.ResponseWriter, request *http.Request) {
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

	response.WriteHeader(status)
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(&signinResponse)
}
