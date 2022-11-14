package signup_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"passwordserver/src/lib/database"
	"passwordserver/src/routes/signup"
	"testing"
)

func TestSignup(t *testing.T) {
	b64MasterHash := base64.StdEncoding.EncodeToString([]byte{0, 0, 0, 0, 0})
	b64ProtectedDatabaseKey := base64.StdEncoding.EncodeToString([]byte{1, 1, 1, 1, 1})

	parameters := signup.SignupParameters{
		Email:                "testEmail",
		MasterHash:           b64MasterHash,
		ProtectedDatabaseKey: b64ProtectedDatabaseKey,
	}
	bufferParameters := new(bytes.Buffer)
	json.NewEncoder(bufferParameters).Encode(parameters)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bufferParameters)
	responseWriter := httptest.NewRecorder()

	signup.SignupHandler(responseWriter, request)

	response := responseWriter.Result()

	defer response.Body.Close()
	body, bodyError := io.ReadAll(response.Body)

	if bodyError != nil {
		t.Errorf("Expected bodyError to be nil, got %v", bodyError)
	}

	parsedErrorBody := signup.SignupErrorResponse{}
	json.NewDecoder(response.Body).Decode(&parsedErrorBody)

	parsedBody := signup.SignupResponse{}
	json.NewDecoder(response.Body).Decode(&parsedBody)

	if parsedErrorBody.Error == "" {
		t.Logf("Signup Success")
	} else {
		t.Errorf("Expected {}, got %v", string(body))
	}
}

func init() {
	os.Setenv("ENVIRONMENT", "testing")
	go database.DatabaseConnect()
}
