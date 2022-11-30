package routes

import (
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"net/http"

	psCrypto "passwordserver/src/lib/crypto"
	psDatabase "passwordserver/src/lib/database"
	psUtils "passwordserver/src/lib/utils"
)

type SigninParameters struct {
	Email      string
	MasterHash string
}

type SigninResponse struct {
	Authenticated bool
}

type SigninErrorResponse struct {
	Error string
}

func SigninPost(response http.ResponseWriter, request *http.Request) {
	signinParameters := SigninParameters{}
	decoderError := json.NewDecoder(request.Body).Decode(&signinParameters)
	if decoderError != nil {
		psUtils.JsonResponse(response, http.StatusBadRequest, SigninErrorResponse{Error: "Unable to decode JSON body."})
		return
	}

	if signinParameters.MasterHash == "" {
		psUtils.JsonResponse(response, http.StatusBadRequest, SigninErrorResponse{Error: "Required parameter 'MasterHash' not provided."})
		return
	}

	if signinParameters.Email == "" {
		psUtils.JsonResponse(response, http.StatusBadRequest, SigninErrorResponse{Error: "Required parameter 'Email' not provided."})
		return
	}

	MasterHashBytes, dmhError := hex.DecodeString(signinParameters.MasterHash)
	if dmhError != nil {
		psUtils.JsonResponse(response, http.StatusBadRequest, SigninErrorResponse{Error: "Unable to decode hex encoded parameter 'MasterHash'."})
		return
	}

	user := psDatabase.User{}
	psDatabase.Database.First(&user, "email = ?", signinParameters.Email)

	strengthenedMasterHashBytes := psCrypto.StrengthenMasterHash(MasterHashBytes, user.MasterHashSalt)
	same := subtle.ConstantTimeCompare(user.MasterHash, strengthenedMasterHashBytes) == 1

	if same {
		cookieError := psCrypto.CreateSessionCookie(response, user)

		if cookieError != nil {
			psUtils.JsonResponse(response, http.StatusInternalServerError, SigninErrorResponse{Error: "Error while setting session cookie."})
			return
		}
	}

	psUtils.JsonResponse(response, http.StatusOK, SigninResponse{Authenticated: same})
}
