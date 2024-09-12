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

	payloadAccess, err := jwtGen.VerifyToken(tokens.Pair.AccessToken, true)
	if err != nil {
		t.Fatalf("expected valid access token, got error %v", err)
	}

	if payloadAccess.TokenUUID != tokens.TokenUUID {
		t.Fatalf("expected access token UUID, got error %v", err)
	}

	if payloadAccess.IPAddress != ip {
		t.Fatalf("expected access token ip adress, got error %v", err)
	}

	if payloadAccess.UserID != userID {
		t.Fatalf("expected access token ip adress, got error %v", err)
	}

	payloadRefresh, err := jwtGen.VerifyToken(tokens.Pair.RefreshToken, false)
	if err != nil {
		t.Fatalf("expected valid refresh token, got error %v", err)
	}

	if payloadRefresh.TokenUUID != tokens.TokenUUID {
		t.Fatalf("expected refresh token UUID, got error %v", err)
	}

	if payloadRefresh.IPAddress != ip {
		t.Fatalf("expected refresh token ip adress, got error %v", err)
	}

	if payloadRefresh.UserID != userID {
		t.Fatalf("expected refresh token ip adress, got error %v", err)
	}
}
