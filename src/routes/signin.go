package routes

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"passwordserver/src/lib"
	libcrypto "passwordserver/src/lib/crypto"
	"passwordserver/src/lib/database"
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

func SigninHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		SigninPost(response, request)
	}
}

func SigninPost(response http.ResponseWriter, request *http.Request) {
	signinParameters := SigninParameters{}
	decoderError := json.NewDecoder(request.Body).Decode(&signinParameters)
	if decoderError != nil {
		lib.JsonResponse(response, http.StatusBadRequest, SigninErrorResponse{Error: "Unable to decode JSON body."})
		return
	}

	if signinParameters.MasterHash == "" {
		lib.JsonResponse(response, http.StatusBadRequest, SigninErrorResponse{Error: "Required parameter 'MasterHash' not provided."})
		return
	}

	if signinParameters.Email == "" {
		lib.JsonResponse(response, http.StatusBadRequest, SigninErrorResponse{Error: "Required parameter 'Email' not provided."})
		return
	}

	MasterHashBytes, dmhError := base64.StdEncoding.DecodeString(signinParameters.MasterHash)
	if dmhError != nil {
		lib.JsonResponse(response, http.StatusBadRequest, SigninErrorResponse{Error: "Unable to decode base64 encoded parameter 'MasterHash'."})
		return
	}

	if database.Database != nil {
		user := database.User{}
		database.Database.Model(database.User{Email: signinParameters.Email}).First(&user)

		strengthenedMasterHashBytes := libcrypto.StrengthenMasterHash(MasterHashBytes, user.MasterHashSalt)
		same := subtle.ConstantTimeCompare(user.MasterHash, strengthenedMasterHashBytes) == 1

		if same {
			libcrypto.CreateSessionCookie(response, user)
		}

		lib.JsonResponse(response, http.StatusOK, SigninResponse{Authenticated: same})
	}
}
