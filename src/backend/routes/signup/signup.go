package signup

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"passwordserver/src/lib"
	libcrypto "passwordserver/src/lib/crypto"
	"passwordserver/src/lib/database"

	"github.com/google/uuid"
)

type SignupParameters struct {
	Email                string
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
	case http.MethodPost:
		SignupPost(response, request)
	}
}

func SignupPost(response http.ResponseWriter, request *http.Request) {
	signupParameters := SignupParameters{}
	decoderError := json.NewDecoder(request.Body).Decode(&signupParameters)
	if decoderError != nil {
		lib.JsonResponse(response, http.StatusBadRequest, SignupErrorResponse{Error: "Unable to decode JSON body."})
		return
	}

	if signupParameters.Email == "" {
		lib.JsonResponse(response, http.StatusBadRequest, SignupErrorResponse{Error: "Required parameter 'Email' not provided."})
		return
	}

	if signupParameters.MasterHash == "" {
		lib.JsonResponse(response, http.StatusBadRequest, SignupErrorResponse{Error: "Required parameter 'MasterHash' not provided."})
		return
	}

	if signupParameters.ProtectedDatabaseKey == "" {
		lib.JsonResponse(response, http.StatusBadRequest, SignupErrorResponse{Error: "Required parameter 'ProtectedDatabaseKey' not provided."})
		return
	}

	MasterHashBytes, dmhError := hex.DecodeString(signupParameters.MasterHash)
	if dmhError != nil {
		lib.JsonResponse(response, http.StatusBadRequest, SignupErrorResponse{Error: "Unable to decode hex encoded parameter 'MasterHash'."})
		return
	}

	strengthenedMasterHashSalt := libcrypto.RandomBytes(16)
	strengthenedMasterHashBytes := libcrypto.StrengthenMasterHash(MasterHashBytes, strengthenedMasterHashSalt)
	decodedProtectedDatabaseKey, _ := hex.DecodeString(signupParameters.ProtectedDatabaseKey)

	if database.Database != nil {
		newUser := database.User{
			Email:                signupParameters.Email,
			MasterHash:           strengthenedMasterHashBytes,
			MasterHashSalt:       strengthenedMasterHashSalt,
			ProtectedDatabaseKey: decodedProtectedDatabaseKey,
		}
		database.Database.Create(&newUser)
		lib.JsonResponse(response, http.StatusOK, SignupResponse{UserId: newUser.Id})
	} else {
		lib.JsonResponse(response, http.StatusInternalServerError, SignupErrorResponse{Error: "The server was unable to create a new user."})
		return
	}
}
