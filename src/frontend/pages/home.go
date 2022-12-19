package pages

import (
	"net/http"

	"passwordserver/src/frontend"
	psCrypto "passwordserver/src/lib/crypto"
	psDatabase "passwordserver/src/lib/database"
	psRedirects "passwordserver/src/lib/redirects"
)

type HomeData struct {
	User psDatabase.User
}

func Home(response http.ResponseWriter, request *http.Request) {
	template := frontend.Template("home.html", "base.html")

	authenticated, user, _, verifyError := psCrypto.VerifySessionCookie(request)

	if verifyError != nil {
		psRedirects.RedirectSignin(response, http.StatusPermanentRedirect)
	}

	if authenticated {
		frontend.ExecuteTemplate(response, template, HomeData{User: user})
	} else {
		psRedirects.RedirectSignin(response, http.StatusPermanentRedirect)
	}
}
