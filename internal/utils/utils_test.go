package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/karelmolina/play5/model"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "secret123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" || hash == password {
		t.Fatal("HashPassword returned empty or plaintext")
	}
	if !CheckPassword(hash, password) {
		t.Fatal("CheckPassword should succeed with correct password")
	}
	if CheckPassword(hash, "wrongpassword") {
		t.Fatal("CheckPassword should fail with wrong password")
	}
}

func TestGenerateAndParseToken(t *testing.T) {
	SetJWTSecret("this-is-a-very-long-secret-key-that-is-32-bytes!")

	user := model.User{
		ID:                uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Username:          "testuser",
		Role:              model.RolePlayer,
		IsApproved:        false,
		PreferredLanguage: "es",
	}

	tokenStr, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("GenerateToken returned empty string")
	}

	claims, err := ParseToken(tokenStr)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	if claims.Sub != user.ID {
		t.Errorf("sub claim mismatch: got %v, want %v", claims.Sub, user.ID)
	}
	if claims.Username != user.Username {
		t.Errorf("username claim mismatch: got %v, want %v", claims.Username, user.Username)
	}
	if claims.Role != string(user.Role) {
		t.Errorf("role claim mismatch: got %v, want %v", claims.Role, user.Role)
	}
	if claims.IsApproved != user.IsApproved {
		t.Errorf("isApproved claim mismatch: got %v, want %v", claims.IsApproved, user.IsApproved)
	}
	if claims.PreferredLanguage != user.PreferredLanguage {
		t.Errorf("preferredLanguage claim mismatch: got %v, want %v", claims.PreferredLanguage, user.PreferredLanguage)
	}

	exp := claims.RegisteredClaims.ExpiresAt
	if exp == nil {
		t.Fatal("exp claim is nil")
	}
	iat := claims.RegisteredClaims.IssuedAt
	if iat == nil {
		t.Fatal("iat claim is nil")
	}

	duration := exp.Time.Sub(iat.Time)
	if duration != 24*time.Hour {
		t.Errorf("token duration mismatch: got %v, want 24h", duration)
	}
}

func TestParseExpiredToken(t *testing.T) {
	SetJWTSecret("this-is-a-very-long-secret-key-that-is-32-bytes!")

	claims := Claims{
		Sub:               uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		Username:          "expired",
		Role:              "player",
		IsApproved:        false,
		PreferredLanguage: "en",
		RegisteredClaims:  jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour))},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(jwtSecret)

	_, err := ParseToken(tokenStr)
	if err == nil {
		t.Fatal("ParseToken should fail for expired token")
	}
}

func TestParseInvalidToken(t *testing.T) {
	SetJWTSecret("this-is-a-very-long-secret-key-that-is-32-bytes!")

	_, err := ParseToken("not.a.token")
	if err == nil {
		t.Fatal("ParseToken should fail for invalid token")
	}
}

func TestValidationStruct(t *testing.T) {
	type TestStruct struct {
		Name string `validate:"required,min=3"`
	}

	valid := TestStruct{Name: "John"}
	if err := ValidateStruct(valid); err != nil {
		t.Errorf("ValidateStruct should pass for valid struct: %v", err)
	}

	invalid := TestStruct{Name: "Jo"}
	if err := ValidateStruct(invalid); err == nil {
		t.Error("ValidateStruct should fail for invalid struct")
	}
}
