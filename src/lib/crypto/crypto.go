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

func CreateSessionCookie(response http.ResponseWriter, user database.User) error {
	if database.Database != nil {
		privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
		publicKey := &privateKey.PublicKey

		sessionToken := database.SessionToken{
			UserId: user.Id,
			N:      publicKey.N.Bytes(),
			E:      publicKey.E,
		}
		database.Database.Create(&sessionToken)
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
			Expires:  time.Now().Add(45784758475874),
			Secure:   false,
			HttpOnly: true,
			Path:     "/",
		}

		http.SetCookie(response, &cookie)

		return nil
	} else {
		return database.Database.Error
	}
}

func VerifySessionCookie(request *http.Request) (bool, database.User) {
	if database.Database != nil {
		cookie, cookieError := request.Cookie("SessionToken")

		if cookieError != nil {
			return false, database.User{}
		}

		splitValue := strings.Split(cookie.Value, ",")
		jsonSessionCookie, _ := hex.DecodeString(splitValue[0])

		signature, _ := hex.DecodeString(splitValue[1])

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

		user := database.User{}
		database.Database.First("id = ?", sessionToken.UserId)

		if rsa.VerifyPKCS1v15(&publicKey, crypto.SHA512, hashed[:], signature) == nil {
			return true, user
		}

		return false, database.User{}
	} else {
		return false, database.User{}
	}
}
