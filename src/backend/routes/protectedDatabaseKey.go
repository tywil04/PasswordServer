package routes

import (
	"encoding/hex"
	"net/http"

	psCrypto "passwordserver/src/lib/crypto"
	psDatabase "passwordserver/src/lib/database"
	psUtils "passwordserver/src/lib/utils"
)

type ProtectedDatabaseKeyResponse struct {
	ProtectedDatabaseKey string
}

type ProtectedDatabaseKeyErrorResponse struct {
	Error string
}

func ProtectedDatabaseKeyGet(response http.ResponseWriter, request *http.Request) {
	if psDatabase.Database != nil {
		authenticated, user, _, _ := psCrypto.VerifySessionCookie(request)

		if authenticated {
			protectedDatabaseKey := hex.EncodeToString(user.ProtectedDatabaseKey)
			psUtils.JsonResponse(response, http.StatusOK, ProtectedDatabaseKeyResponse{ProtectedDatabaseKey: protectedDatabaseKey})
			return
		} else {
			psUtils.JsonResponse(response, http.StatusNetworkAuthenticationRequired, ProtectedDatabaseKeyErrorResponse{Error: "You are not authenticated."})
			return
		}
	} else {
		psUtils.JsonResponse(response, http.StatusInternalServerError, ProtectedDatabaseKeyErrorResponse{Error: "The server was unable to retrieve user."})
		return
	}
}
