package auth

import (
	"testing"

	"github.com/cmackin9500/goecom/config"
)

func TestCreateJWT(t *testing.T) {
	secret := []byte(config.Envs.JWTSecret)
	token, err := CreateJWT(secret, 1)

	if err != nil {
		t.Errorf("failed to create JWT: %v", err)
	}

	if token == "" {
		t.Errorf("created blank JWT: %v", err)
	}
}