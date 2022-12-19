package routes

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	psCrypto "passwordserver/src/lib/crypto"
	psDatabase "passwordserver/src/lib/database"
	psUtils "passwordserver/src/lib/utils"

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

func SignupPost(response http.ResponseWriter, request *http.Request) {
	signupParameters := SignupParameters{}
	decoderError := json.NewDecoder(request.Body).Decode(&signupParameters)
	if decoderError != nil {
		psUtils.JsonResponse(response, http.StatusBadRequest, SignupErrorResponse{Error: "Unable to decode JSON body."})
		return
	}

	if signupParameters.Email == "" {
		psUtils.JsonResponse(response, http.StatusBadRequest, SignupErrorResponse{Error: "Required parameter 'Email' not provided."})
		return
	}

	if signupParameters.MasterHash == "" {
		psUtils.JsonResponse(response, http.StatusBadRequest, SignupErrorResponse{Error: "Required parameter 'MasterHash' not provided."})
		return
	}

	if signupParameters.ProtectedDatabaseKey == "" {
		psUtils.JsonResponse(response, http.StatusBadRequest, SignupErrorResponse{Error: "Required parameter 'ProtectedDatabaseKey' not provided."})
		return
	}

	MasterHashBytes, dmhError := hex.DecodeString(signupParameters.MasterHash)
	if dmhError != nil {
		psUtils.JsonResponse(response, http.StatusBadRequest, SignupErrorResponse{Error: "Unable to decode hex encoded parameter 'MasterHash'."})
		return
	}

	strengthenedMasterHashSalt := psCrypto.RandomBytes(16)
	strengthenedMasterHashBytes := psCrypto.StrengthenMasterHash(MasterHashBytes, strengthenedMasterHashSalt)

	protectedDatabaseParts := strings.Split(signupParameters.ProtectedDatabaseKey, ";")
	protectedDatabaseKey := protectedDatabaseParts[1]
	protectedDatabaseKeyIV := protectedDatabaseParts[0]

	decodedProtectedDatabaseKey, _ := hex.DecodeString(protectedDatabaseKey)
	decodedProtectedDatabaseKeyIV, _ := hex.DecodeString(protectedDatabaseKeyIV)

	newUser := psDatabase.User{
		Email:                  signupParameters.Email,
		MasterHash:             strengthenedMasterHashBytes,
		MasterHashSalt:         strengthenedMasterHashSalt,
		ProtectedDatabaseKey:   decodedProtectedDatabaseKey,
		ProtectedDatabaseKeyIV: decodedProtectedDatabaseKeyIV,
	}
	psDatabase.Database.Create(&newUser)
	psUtils.JsonResponse(response, http.StatusOK, SignupResponse{UserId: newUser.Id})
}
