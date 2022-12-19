package signin_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"passwordserver/src/backend/routes"

	psDatabase "passwordserver/src/lib/database"
)

func TestSigninWithGoodCredentials(t *testing.T) {
	// these values have been pre-generated against a known-good client
	// email: testEmail@testEmail.co.uk
	// password: testPassword
	email := "testEmail@testEmail.co.uk"
	masterHash := "d7a292c897c0dd1e31daa7338bfed97e5036a60d30ca1bd8c8b4d05b778ae3e3b82ef1d393b1ac368cde62f8041d204ef8f620b6107e9f0d7e1f3ffef8f95280"

	strengthenedMasterHash := "62b961c21e687740953bb772aa5d823cd03a13c8075e7f067e7c9188b8dbd97237b47c13222a6e6164c2953e10e05bf18b32b6946344ffd8efcda6eaaea69a6d"
	strengthenedMasterHashSalt := "727e92b72d25c31735a18c2cc3126fd0"
	protectedDatabaseKey := "13407cbaf62d537a4bcf2f7b69a745a12d3873014f21c54bb43b4b4c8a82f8d3eff37e6c67d83f8abbb1525445d5b9d6"
	protectedDatabaseKeyIV := "ce0db8e578c2fa845570e2c366df8707"

	decodedStrengthenedMasterHash, _ := hex.DecodeString(strengthenedMasterHash)
	decodedStrengthenedMasterHashSalt, _ := hex.DecodeString(strengthenedMasterHashSalt)
	decodedProtectedDatabaseKey, _ := hex.DecodeString(protectedDatabaseKey)
	decodedProtectedDatabaseKeyIV, _ := hex.DecodeString(protectedDatabaseKeyIV)

	newUser := psDatabase.User{
		Email:                  email,
		MasterHash:             decodedStrengthenedMasterHash,
		MasterHashSalt:         decodedStrengthenedMasterHashSalt,
		ProtectedDatabaseKey:   decodedProtectedDatabaseKey,
		ProtectedDatabaseKeyIV: decodedProtectedDatabaseKeyIV,
	}
	psDatabase.Database.Create(&newUser)

	// start testing
	parameters := routes.SignupParameters{
		Email:                email,
		MasterHash:           masterHash,
		ProtectedDatabaseKey: protectedDatabaseKeyIV + ";" + protectedDatabaseKey,
	}
	bufferParameters := new(bytes.Buffer)
	json.NewEncoder(bufferParameters).Encode(parameters)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signin", bufferParameters)
	responseWriter := httptest.NewRecorder()

	routes.SigninPost(responseWriter, request)

	response := responseWriter.Result()

	body, bodyError := io.ReadAll(response.Body)
	defer response.Body.Close()

	if bodyError != nil {
		t.Errorf("Expected bodyError to be nil, got %v", bodyError)
	}

	decoded := map[string]any{}
	json.NewDecoder(bytes.NewBuffer(body)).Decode(&decoded)

	if decoded["Authenticated"] == true {
		t.Logf("Signup Success")
	} else {
		t.Errorf(`Expected {"Authenticated":true}, got %v`, string(body))
	}
}

func init() {
	os.Setenv("ENVIRONMENT", "testing")
	psDatabase.DatabaseConnect()
}
