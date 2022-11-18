package routes

import (
	"net/http"

	"passwordserver/src/lib"
	libcrypto "passwordserver/src/lib/crypto"
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
	signedOut, cscError := libcrypto.ClearSessionCookie(response, request)

	if cscError != nil {
		lib.JsonResponse(response, http.StatusBadRequest, SignoutErrorResponse{Error: "Error while trying to clear session cookie."})
		return
	}

	lib.JsonResponse(response, http.StatusOK, SignoutResponse{SignedOut: signedOut})
}
