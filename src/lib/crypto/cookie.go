package crypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	psDatabase "passwordserver/src/lib/database"

	"github.com/google/uuid"
)

type SessionCookie struct {
	UserId         uuid.UUID
	SessionTokenId uuid.UUID
}

func CreateSessionCookie(response http.ResponseWriter, user psDatabase.User) error {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	publicKey := &privateKey.PublicKey

	cookieExpiry := time.Now().Add(time.Hour)

	sessionToken := psDatabase.SessionToken{
		UserId: user.Id,
		N:      publicKey.N.Bytes(),
		E:      publicKey.E,
		Expiry: cookieExpiry,
	}
	psDatabase.Database.Create(&sessionToken)
	user.SessionTokens = append(user.SessionTokens, sessionToken)

	sessionCookie := SessionCookie{
		UserId:         user.Id,
		SessionTokenId: sessionToken.Id,
	}
	jsonPayload := new(bytes.Buffer)
	json.NewEncoder(jsonPayload).Encode(sessionCookie)
	hashed := sha512.Sum512(jsonPayload.Bytes())

	signature, _ := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA512, hashed[:])

	encodedSignature := hex.EncodeToString(signature)
	encodedSessionCookie := hex.EncodeToString(jsonPayload.Bytes())

	cookie := http.Cookie{
		Name:     "SessionToken",
		Value:    encodedSessionCookie + "," + encodedSignature,
		Expires:  cookieExpiry,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(response, &cookie)

	return nil
}

func VerifySessionCookie(request *http.Request) (bool, psDatabase.User, psDatabase.SessionToken, error) {
	cookie, cookieError := request.Cookie("SessionToken")

	if cookie == nil || cookie.Value == "" {
		return false, psDatabase.User{}, psDatabase.SessionToken{}, nil
	}

	if cookieError != nil {
		return false, psDatabase.User{}, psDatabase.SessionToken{}, cookieError
	}

	splitValue := strings.Split(cookie.Value, ",")
	jsonSessionCookie, _ := hex.DecodeString(splitValue[0])

	signature, _ := hex.DecodeString(splitValue[1])

	sessionCookie := SessionCookie{}
	json.NewDecoder(bytes.NewBuffer(jsonSessionCookie)).Decode(&sessionCookie)

	sessionToken := psDatabase.SessionToken{}
	psDatabase.Database.First(&sessionToken, "id = ?", sessionCookie.SessionTokenId, "user_id = ?", sessionCookie.UserId)

	publicKey := rsa.PublicKey{
		N: new(big.Int).SetBytes(sessionToken.N),
		E: sessionToken.E,
	}

	jsonPayload := new(bytes.Buffer)
	json.NewEncoder(jsonPayload).Encode(sessionCookie)
	hashed := sha512.Sum512(jsonPayload.Bytes())

	allSessionTokens := []psDatabase.SessionToken{}
	psDatabase.Database.Find(&allSessionTokens)
	for _, sessionToken := range allSessionTokens {
		// if the session has expired remove the token. we are lax about this because any cookie which is expired doesn't isnt valid which means the code cant even process it
		if sessionToken.Expiry.Before(time.Now()) {
			psDatabase.Database.Delete(&sessionToken)
		}
	}

	if rsa.VerifyPKCS1v15(&publicKey, crypto.SHA512, hashed[:], signature) == nil {
		user := psDatabase.User{}
		psDatabase.Database.First(&user, "id = ?", sessionToken.UserId)

		return true, user, sessionToken, nil
	}

	return false, psDatabase.User{}, psDatabase.SessionToken{}, nil
}

func ClearSessionCookie(response http.ResponseWriter, request *http.Request) (bool, error) {
	authenticated, _, sessionToken, _ := VerifySessionCookie(request)

	if authenticated {
		cookie := http.Cookie{
			Name:     "SessionToken",
			Value:    "",
			MaxAge:   -1,
			Secure:   os.Getenv("ENVIRONMENT") == "production",
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
		}

		http.SetCookie(response, &cookie)

		psDatabase.Database.Delete(&sessionToken)

		return true, nil
	}

	return false, nil
}
