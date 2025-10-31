package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	accessTokenTTL  = 60 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

var jwtKey = []byte("SecretKey")

type JwtClaim struct {
	UserID uuid.UUID `json:"user_id"`
	Login  string    `json:"login"`
	Type   string    `json:"type"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uuid.UUID, login string, tokenType string) (string, error) {
	var expiredTime time.Time
	if tokenType == "access" {
		expiredTime = time.Now().Add(accessTokenTTL)
	} else {
		expiredTime = time.Now().Add(refreshTokenTTL)
	}

	claims := &JwtClaim{
		UserID: userID,
		Login:  login,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*JwtClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
