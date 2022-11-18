package routes

import (
	"net/http"

	psCrypto "passwordserver/src/lib/crypto"
	psUtils "passwordserver/src/lib/utils"
)

type SignoutResponse struct {
	SignedOut bool
}

type SignoutErrorResponse struct {
	Error string
}

func SignoutHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodDelete:
		SignoutDelete(response, request)
	}
}

func SignoutDelete(response http.ResponseWriter, request *http.Request) {
	signedOut, cscError := psCrypto.ClearSessionCookie(response, request)

	if cscError != nil {
		psUtils.JsonResponse(response, http.StatusBadRequest, SignoutErrorResponse{Error: "Error while trying to clear session cookie."})
		return
	}

	psUtils.JsonResponse(response, http.StatusOK, SignoutResponse{SignedOut: signedOut})
}
