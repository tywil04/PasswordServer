package crypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"net/http"
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

func GenerateSessionCookie(privateKey rsa.PrivateKey, sessionTokenId uuid.UUID, userId uuid.UUID) http.Cookie {
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
		Expires:  time.Now().Add(365),
		Secure:   false,
		HttpOnly: true,
		Path:     "/",
	}

	return cookie
}
