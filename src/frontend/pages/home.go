package pages

import (
	"io/fs"
	"net/http"
	"text/template"

	psCrypto "passwordserver/src/lib/crypto"
	psDatabase "passwordserver/src/lib/database"
	psRedirect "passwordserver/src/lib/redirect"
)

type HomeData struct {
	User psDatabase.User
}

func Home(response http.ResponseWriter, request *http.Request, htmlDir fs.FS) {
	template := template.Must(template.ParseFS(htmlDir, "home.html", "base.html"))

	authenticated, user, _, verifyError := psCrypto.VerifySessionCookie(request)

	if verifyError != nil {
		psRedirect.RedirectSignin(response, http.StatusPermanentRedirect)
	}

	if authenticated {
		template.Execute(response, HomeData{User: user})
	} else {
		psRedirect.RedirectSignin(response, http.StatusPermanentRedirect)
	}
}
