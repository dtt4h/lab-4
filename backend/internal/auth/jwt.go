package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func NewAccessToken(userID, secret string, ttl time.Duration) (string, error) {
	claims := Claims{TokenType: "access", RegisteredClaims: jwt.RegisteredClaims{Subject: userID, IssuedAt: jwt.NewNumericDate(time.Now()), ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl))}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func ParseAccessToken(tokenString, secret string) (string, error) {
	return parseToken(tokenString, secret, "access")
}

func NewRefreshToken(userID, secret string, ttl time.Duration) (string, error) {
	claims := Claims{TokenType: "refresh", RegisteredClaims: jwt.RegisteredClaims{Subject: userID, IssuedAt: jwt.NewNumericDate(time.Now()), ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl))}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func ParseRefreshToken(tokenString, secret string) (string, error) {
	return parseToken(tokenString, secret, "refresh")
}

func parseToken(tokenString, secret, tokenType string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.TokenType != tokenType || claims.Subject == "" {
		return "", fmt.Errorf("invalid access token")
	}
	return claims.Subject, nil
}
