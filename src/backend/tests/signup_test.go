package signup_test

import (
	"os"
	"testing"

	psDatabase "passwordserver/src/lib/database"
)

func TestSignup(t *testing.T) {

}

func init() {
	os.Setenv("ENVIRONMENT", "testing")
	go psDatabase.DatabaseConnect()
}
