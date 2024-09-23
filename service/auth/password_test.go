package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error while hashing password: %v", err)
	}

	if hash  == "" {
		t.Errorf("hashed password is blank")
	}

	if hash == "password" {
		t.Errorf("hash and password are the same")
	}
}

func TestComparePassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error while hashing password: %v", err)
	}

	if !ComparePassowords(hash, []byte("password")) {
		t.Errorf("password is not matching even though it should")
	}

	if ComparePassowords(hash, []byte("wrongpassword")) {
		t.Errorf("password matching even though it shouldn't")
	}
}