package routes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"passwordserver/src/lib"
)

type SignupParameters struct {
	MasterHash           string
	ProtectedDatabaseKey string
}

type SignupResponse struct {
}

type SignupErrorResponse struct {
	Error string
}

func SignupHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		SignupPost(response, request)
	}
}

func SignupPost(response http.ResponseWriter, request *http.Request) {
	SignupParameters := SignupParameters{}
	decoderError := json.NewDecoder(request.Body).Decode(&SignupParameters)
	if decoderError != nil {
		lib.JsonResponse(response, http.StatusBadRequest, SignupErrorResponse{Error: "Unable to decode JSON body."})
		return
	}

	if SignupParameters.MasterHash == "" {
		lib.JsonResponse(response, http.StatusBadRequest, SignupErrorResponse{Error: "Required parameter 'MasterHash' not provided."})
		return
	}

	if SignupParameters.ProtectedDatabaseKey == "" {
		lib.JsonResponse(response, http.StatusBadRequest, SignupErrorResponse{Error: "Required parameter 'ProtectedDatabaseKey' not provided."})
		return
	}

	MasterHashBytes, dmhError := base64.StdEncoding.DecodeString(SignupParameters.MasterHash)
	if dmhError != nil {
		lib.JsonResponse(response, http.StatusBadRequest, SignupErrorResponse{Error: "Unable to decode base64 encoded parameter 'MasterHash'."})
		return
	}

	strengthenedMasterHashBytes, strengthenedMasterHashSalt := lib.StrengthenMasterHash(MasterHashBytes)
	strengthenedMasterHash := base64.StdEncoding.EncodeToString(strengthenedMasterHashBytes) + ";" + base64.StdEncoding.EncodeToString(strengthenedMasterHashSalt)

	fmt.Println(strengthenedMasterHash)
	fmt.Println(SignupParameters.ProtectedDatabaseKey)

	lib.JsonResponse(response, http.StatusOK, SignupResponse{})
}
