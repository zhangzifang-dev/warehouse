package jwt

import (
	"testing"
	"time"
)

func TestNewJWT(t *testing.T) {
	secret := "test-secret-key"
	expire := time.Hour

	j := NewJWT(secret, expire)
	if j == nil {
		t.Fatal("NewJWT returned nil")
	}
	if string(j.secret) != secret {
		t.Errorf("expected secret %s, got %s", secret, string(j.secret))
	}
	if j.expire != expire {
		t.Errorf("expected expire %v, got %v", expire, j.expire)
	}
}

func TestGenerateAndParseToken(t *testing.T) {
	j := NewJWT("test-secret-key", time.Hour)

	userID := int64(123)
	username := "testuser"

	token, err := j.GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken returned empty token")
	}

	claims, err := j.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected UserID %d, got %d", userID, claims.UserID)
	}
	if claims.Username != username {
		t.Errorf("expected Username %s, got %s", username, claims.Username)
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	j := NewJWT("test-secret-key", time.Hour)

	_, err := j.ParseToken("invalid-token")
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	j1 := NewJWT("secret1", time.Hour)
	j2 := NewJWT("secret2", time.Hour)

	token, err := j1.GenerateToken(1, "user")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = j2.ParseToken(token)
	if err == nil {
		t.Error("expected error when parsing token with wrong secret, got nil")
	}
}

func TestParseToken_ExpiredToken(t *testing.T) {
	j := NewJWT("test-secret-key", -time.Hour)

	token, err := j.GenerateToken(1, "user")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	time.Sleep(time.Second)

	_, err = j.ParseToken(token)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestParseToken_EmptyToken(t *testing.T) {
	j := NewJWT("test-secret-key", time.Hour)

	_, err := j.ParseToken("")
	if err == nil {
		t.Error("expected error for empty token, got nil")
	}
}

func TestGenerateToken_DifferentUsers(t *testing.T) {
	j := NewJWT("test-secret-key", time.Hour)

	token1, err := j.GenerateToken(1, "user1")
	if err != nil {
		t.Fatalf("GenerateToken failed for user1: %v", err)
	}

	token2, err := j.GenerateToken(2, "user2")
	if err != nil {
		t.Fatalf("GenerateToken failed for user2: %v", err)
	}

	if token1 == token2 {
		t.Error("tokens for different users should be different")
	}
}
