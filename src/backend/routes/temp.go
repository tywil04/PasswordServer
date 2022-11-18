package routes

import (
	"fmt"
	"net/http"

	psCrypto "passwordserver/src/lib/crypto"
)

func TempGet(response http.ResponseWriter, request *http.Request) {
	authenticated, user, _, _ := psCrypto.VerifySessionCookie(request)

	if authenticated {
		fmt.Println(user.Email)
	} else {
		fmt.Println("User not authenticated")
	}
}
