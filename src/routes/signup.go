package routes

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"passwordserver/src/lib"

	"github.com/google/uuid"
)

type SignupParameters struct {
	MasterHash           string
	ProtectedDatabaseKey string
}

type SignupResponse struct {
	UserId uuid.UUID
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
	decodedProtectedDatabaseKey, _ := base64.StdEncoding.DecodeString(SignupParameters.ProtectedDatabaseKey)

	if lib.Database != nil {
		newUser := lib.User{
			MasterHash:           strengthenedMasterHashBytes,
			MasterHashSalt:       strengthenedMasterHashSalt,
			ProtectedDatabaseKey: decodedProtectedDatabaseKey,
		}
		lib.Database.Create(&newUser)
		lib.JsonResponse(response, http.StatusOK, SignupResponse{UserId: newUser.Id})
	} else {
		lib.JsonResponse(response, http.StatusInternalServerError, SignupErrorResponse{Error: "The server was unable to create a new user."})
		return
	}
}
