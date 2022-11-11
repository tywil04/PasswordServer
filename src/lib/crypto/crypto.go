package crypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"passwordserver/src/lib/database"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
)

func StrengthenMasterHash(masterHash []byte, salt []byte) []byte {
	return pbkdf2.Key(masterHash, salt, 150000, 512/8, sha512.New)
}

func RandomBytes(byteLength int) []byte {
	bytes := make([]byte, byteLength)
	rand.Read(bytes)
	return bytes
}

type SessionCookie struct {
	UserId         uuid.UUID
	SessionTokenId uuid.UUID
}

func CreateSessionCookie(response http.ResponseWriter, privateKey rsa.PrivateKey, sessionTokenId uuid.UUID, userId uuid.UUID) {
	sessionCookie := SessionCookie{
		UserId:         userId,
		SessionTokenId: sessionTokenId,
	}
	jsonPayload := new(bytes.Buffer)
	json.NewEncoder(jsonPayload).Encode(sessionCookie)
	hashed := sha512.Sum512(jsonPayload.Bytes())

	signature, _ := rsa.SignPKCS1v15(rand.Reader, &privateKey, crypto.SHA512, hashed[:])

	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	encodedSessionCookie := base64.StdEncoding.EncodeToString(jsonPayload.Bytes())

	cookie := http.Cookie{
		Name:     "SessionToken",
		Value:    encodedSessionCookie + "," + encodedSignature,
		Expires:  time.Now().Add(45784758475874),
		Secure:   false,
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(response, &cookie)
}

func VerifySessionCookie(request *http.Request) bool {
	cookie, cookieError := request.Cookie("SessionToken")

	if cookieError != nil {
		return false
	}

	splitValue := strings.Split(cookie.Value, ",")
	jsonSessionCookie, _ := base64.StdEncoding.DecodeString(splitValue[0])

	signature, _ := base64.StdEncoding.DecodeString(splitValue[1])

	sessionCookie := SessionCookie{}
	json.NewDecoder(bytes.NewBuffer(jsonSessionCookie)).Decode(&sessionCookie)

	sessionToken := database.SessionToken{}
	database.Database.First(&sessionToken, "id = ?", sessionCookie.SessionTokenId, "user_id = ?", sessionCookie.UserId)

	publicKey := rsa.PublicKey{
		N: new(big.Int).SetBytes(sessionToken.N),
		E: sessionToken.E,
	}

	jsonPayload := new(bytes.Buffer)
	json.NewEncoder(jsonPayload).Encode(sessionCookie)
	hashed := sha512.Sum512(jsonPayload.Bytes())

	return rsa.VerifyPKCS1v15(&publicKey, crypto.SHA512, hashed[:], signature) == nil
}
