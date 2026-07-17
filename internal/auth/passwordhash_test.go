package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("securepassword")
	if err != nil {
		t.Fatalf("test failed: HashPassword function error: %v", err)
	}

	if hash == "" {
		t.Fatalf("test failed: returned empty hash string")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "securepassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatal("error making hash, test could not be run")
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("test failed: CheckPasswordHash function error: %v", err)
	}
	if !match {
		t.Fatalf("test failed: password does not match it's own hash")
	}
}

func TestCheckWrongPasswordHash(t *testing.T) {
	password := "securepassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatal("error making hash, test could not be run")
	}

	match, err := CheckPasswordHash("unsafepassword", hash)
	if err != nil {
		t.Fatalf("test failed: CheckPasswordHash function error: %v", err)
	}
	if match {
		t.Fatalf("test failed: hash matches incorrect password")
	}
}
