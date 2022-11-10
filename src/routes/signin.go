package routes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"passwordserver/src/lib"
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

	MasterHashBytes, dmhError := base64.StdEncoding.DecodeString(signinParameters.MasterHash)
	if dmhError != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(SigninErrorResponse{Error: "Unable to decode base64 encoded parameter 'MasterHash'."})
		return
	}

	strengthenedMasterHashBytes, strengthenedMasterHashSalt := lib.StrengthenMasterHash(MasterHashBytes)
	strengthenedMasterHash := base64.StdEncoding.EncodeToString(strengthenedMasterHashBytes) + ";" + base64.StdEncoding.EncodeToString(strengthenedMasterHashSalt)

	fmt.Print(strengthenedMasterHash)
	fmt.Print(signinParameters.ProtectedDatabaseKey)

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(SigninResponse{})
}
