package utils

import (
	"testing"
	"time"
)

// TestGenerateToken generates a new token
func TestGenerateToken(t *testing.T) {
	secret := "secret"
	expiry := 1 * time.Hour
	token, err := GenerateToken(1, "admin", secret, expiry)
	if err != nil {
		t.Errorf("failed to generate token: %v", err)
	}
	if token == "" {
		t.Errorf("token is empty")
	}
}

// TestParseToken parses a token
func TestParseToken(t *testing.T) {
	secret := "secret"
	expiry := 1 * time.Hour
	token, err := GenerateToken(1, "admin", secret, expiry)
	if err != nil {
		t.Errorf("failed to generate token: %v", err)
	}
	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Errorf("failed to parse token: %v", err)
	}
	if claims.UserID != 1 {
		t.Errorf("unexpected user ID: %d", claims.UserID)
	}
	if claims.Role != "admin" {
		t.Errorf("unexpected role: %s", claims.Role)
	}
}

// TestInvalidToken tests with an invalid token
func TestInvalidToken(t *testing.T) {
	secret := "secret"
	_, err := ParseToken("invalid-token", secret)
	if err == nil {
		t.Errorf("expected error for invalid token")
	}
}
