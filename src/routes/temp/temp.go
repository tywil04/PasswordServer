package temp

import (
	"fmt"
	"net/http"
	libcrypto "passwordserver/src/lib/crypto"
)

func SignupHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		TempPost(response, request)
	}
}

func TempPost(response http.ResponseWriter, request *http.Request) {
	authenticated, user := libcrypto.VerifySessionCookie(request)

	if authenticated {
		fmt.Println(user.Email)
	}
}
