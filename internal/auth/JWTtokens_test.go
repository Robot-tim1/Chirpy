package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	mockid := uuid.New()
	tokenSecret := "supersecretcode"

	_, err := MakeJWT(mockid, tokenSecret, time.Second*5)
	if err != nil {
		t.Fatalf("test failed: error making JWT token: %v", err)
	}
}

func TestValidateJWT(t *testing.T) {
	mockid := uuid.New()
	tokenSecret := "supersecretcode"

	tokenString, err := MakeJWT(mockid, tokenSecret, time.Second*5)
	if err != nil {
		t.Fatalf("error making JWT token: %v", err)
	}

	returnid, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("test failed: error validating JWT string %v", err)
	}

	if returnid != mockid {
		t.Fatal("test failed: returned ID does not match input ID")
	}
}

func TestJWTTokenExpire(t *testing.T) {
	mockid := uuid.New()
	tokenSecret := "supersecretcode"

	tokenString, err := MakeJWT(mockid, tokenSecret, time.Millisecond*5)
	if err != nil {
		t.Fatalf("error making JWT token: %v", err)
	}

	time.Sleep(time.Millisecond * 5)

	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Fatal("test failed: token validated when expired")
	}
}

func TestJWTWrongSecret(t *testing.T) {
	mockid := uuid.New()
	tokenSecret := "supersecretcode"

	tokenString, err := MakeJWT(mockid, tokenSecret, time.Second*5)
	if err != nil {
		t.Fatalf("error making JWT token: %v", err)
	}

	_, err = ValidateJWT(tokenString, "superwrongcode")
	if err == nil {
		t.Fatal("test failed: token validated with wrong tokensecret")
	}
}

func TestGetBearerToken(t *testing.T) {
	var value []string
	value = append(value, "Bearer 12345")
	headers := http.Header{"Authorization": value}

	tokenString, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("test failed: error getting token string from header: %v", err)
	}

	if tokenString != "12345" {
		t.Fatalf("test failed: token string output: %s, expected: 12345", tokenString)
	}
}

func TestGetBearerTokenError(t *testing.T) {
	var value []string
	value = append(value, "Bearer 12345")
	headers := http.Header{"Accept": value}

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatalf("test failed: returned no errors with no Authorization header")
	}
}
