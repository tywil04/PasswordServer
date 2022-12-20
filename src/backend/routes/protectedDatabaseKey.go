package routes

import (
	"encoding/hex"
	"net/http"

	psCrypto "passwordserver/src/lib/crypto"
	psUtils "passwordserver/src/lib/utils"
)

type ProtectedDatabaseKeyResponse struct {
	ProtectedDatabaseKey string
}

type ProtectedDatabaseKeyErrorResponse struct {
	Error string
}

func ProtectedDatabaseKeyGet(response http.ResponseWriter, request *http.Request) {
	authenticated, user, _, _ := psCrypto.VerifySessionCookie(request)

	if authenticated {
		protectedDatabaseKey := hex.EncodeToString(user.ProtectedDatabaseKey)
		protectedDatabaseKeyIV := hex.EncodeToString(user.ProtectedDatabaseKeyIV)
		psUtils.JsonResponse(response, http.StatusOK, ProtectedDatabaseKeyResponse{ProtectedDatabaseKey: protectedDatabaseKeyIV + ";" + protectedDatabaseKey})
		return
	} else {
		psUtils.JsonResponse(response, http.StatusNetworkAuthenticationRequired, ProtectedDatabaseKeyErrorResponse{Error: "You are not authenticated."})
		return
	}
}
