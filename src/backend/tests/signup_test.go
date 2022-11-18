package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"passwordserver/src/backend/routes"

	psDatabase "passwordserver/src/lib/database"
)

func TestSignup(t *testing.T) {
	b64MasterHash := "bb01afb3beea90574f8397153ff11b38f80e4f64635c2a21f665a903ce18cdf4f86dff9157ac40e88e052beb3107118a442f37be192162225d8b66b1b4df2982"
	b64ProtectedDatabaseKey := "bc7decd789599be8b7f11b1f87c1c6d6;c098ab38579f9c4f4e091361de8273ef0e7f27f5a665c9248baea1554cf36682466cbb5eb2dc97fd5dbf5cce40c0f849"

	parameters := routes.SignupParameters{
		Email:                "testEmail",
		MasterHash:           b64MasterHash,
		ProtectedDatabaseKey: b64ProtectedDatabaseKey,
	}
	bufferParameters := new(bytes.Buffer)
	json.NewEncoder(bufferParameters).Encode(parameters)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bufferParameters)
	responseWriter := httptest.NewRecorder()

	routes.SignupHandler(responseWriter, request)

	response := responseWriter.Result()

	defer response.Body.Close()
	body, bodyError := io.ReadAll(response.Body)

	if bodyError != nil {
		t.Errorf("Expected bodyError to be nil, got %v", bodyError)
	}

	parsedErrorBody := routes.SignupErrorResponse{}
	json.NewDecoder(response.Body).Decode(&parsedErrorBody)

	parsedBody := routes.SignupResponse{}
	json.NewDecoder(response.Body).Decode(&parsedBody)

	if parsedErrorBody.Error == "" {
		t.Logf("Signup Success")
	} else {
		t.Errorf("Expected {}, got %v", string(body))
	}
}

func init() {
	os.Setenv("ENVIRONMENT", "testing")
	go psDatabase.DatabaseConnect()
}
