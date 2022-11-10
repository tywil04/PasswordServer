package routes

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type SigninParameters struct {
	MasterHash           string
	ProtectedDatabaseKey string
}

type SigninResponse struct {
}

type SigninErrorResponse struct {
	Error string
}

func SigninHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		SigninPost(response, request)
	}
}

func SigninPost(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	signinParameters := SigninParameters{}
	decoderError := json.NewDecoder(request.Body).Decode(&signinParameters)
	if decoderError != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(SigninErrorResponse{Error: "Unable to decode JSON body."})
		return
	}

	if signinParameters.MasterHash == "" {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(SigninErrorResponse{Error: "Required parameter 'MasterHash' not provided."})
		return
	}

	if signinParameters.ProtectedDatabaseKey == "" {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(SigninErrorResponse{Error: "Required parameter 'ProtectedDatabaseKey' not provided."})
		return
	}

	decodedMasterHash, dmhError := base64.StdEncoding.DecodeString(signinParameters.MasterHash)
	decodedProtectedDatabaseKey, dpdkError := base64.StdEncoding.DecodeString(signinParameters.ProtectedDatabaseKey)

	if dmhError != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(SigninErrorResponse{Error: "Unable to decode base64 encoded parameter 'MasterHash'."})
		return
	}

	if dpdkError != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(SigninErrorResponse{Error: "Unable to decode base64 encoded parameter 'ProtectedDatabaseKey'."})
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(SigninResponse{})
}
