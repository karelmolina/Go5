package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/karelmolina/play5/model"
)

var jwtSecret []byte

func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

type Claims struct {
	Sub               uuid.UUID `json:"sub"`
	Username          string    `json:"username"`
	Role              string    `json:"role"`
	IsApproved        bool      `json:"isApproved"`
	PreferredLanguage string    `json:"preferredLanguage"`
	jwt.RegisteredClaims
}

func GenerateToken(user model.User) (string, error) {
	now := time.Now()
	claims := Claims{
		Sub:               user.ID,
		Username:          user.Username,
		Role:              string(user.Role),
		IsApproved:        user.IsApproved,
		PreferredLanguage: user.PreferredLanguage,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token claims")
}
