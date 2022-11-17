package temp

import (
	"fmt"
	"net/http"
	libcrypto "passwordserver/src/lib/crypto"
)

func TempHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		TempPost(response, request)
	}
}

func TempPost(response http.ResponseWriter, request *http.Request) {
	authenticated, user := libcrypto.VerifySessionCookie(request)

	if authenticated {
		fmt.Println(user.Email)
	} else {
		fmt.Println("User not authenticated")
	}
}
