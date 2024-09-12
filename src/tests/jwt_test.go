package tests

import (
	"backend/src/pkg/jwt"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateToken(t *testing.T) {
	var jwtGen jwt.JWTGenerator

	jwtGen.Init("mysecret")

	userID := uuid.New()
	ip := "127.0.0.1"

	tokens, err := jwtGen.CreateTokenPair(userID, ip)

	if err != nil {
		t.Fatalf("expected valid token, got error %v", err)
	}

	payload, err := jwtGen.VerifyToken(tokens.Pair.AccessToken, true)
	if err != nil {
		t.Fatalf("expected valid token, got error %v", err)
	}

	if payload.IPAddress != ip {
		t.Fatalf("expected ip adress, got error %v", err)
	}

	if payload.UserID != userID {
		t.Fatalf("expected ip adress, got error %v", err)
	}
}
