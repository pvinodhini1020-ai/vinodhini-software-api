package tests

import (
	"testing"

	"github.com/vinodhini/software-api/pkg/utils"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	if hash == password {
		t.Error("Hash should not equal plain password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"
	hash, _ := utils.HashPassword(password)

	if !utils.CheckPassword(password, hash) {
		t.Error("Password check should return true for correct password")
	}

	if utils.CheckPassword("wrongpassword", hash) {
		t.Error("Password check should return false for incorrect password")
	}
}

func TestGenerateToken(t *testing.T) {
	token, err := utils.GenerateToken(1, "test@example.com", "admin", "secret", 3600)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}
}

func TestValidateToken(t *testing.T) {
	secret := "test-secret"
	token, _ := utils.GenerateToken(1, "test@example.com", "admin", secret, 3600)

	claims, err := utils.ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", claims.UserID)
	}

	if claims.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", claims.Email)
	}

	if claims.Role != "admin" {
		t.Errorf("Expected role admin, got %s", claims.Role)
	}
}
