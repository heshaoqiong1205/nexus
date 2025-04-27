package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte
var appName string

func init() {
	secretKey = []byte("nexus_secret_key")
	appName = "nexus"
}

type NexcusHeader struct {
	Type string `json:"typ"`
	Alg  string `json:"alg"`
}

type NexusClaims struct {
	NexcusHeader
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey []byte
}

func GenerateToken(algorithm string, userID string, lifetime int64) (string, error) {
	claims := NexusClaims{
		NexcusHeader: NexcusHeader{
			Type: "JWT",
			Alg:  algorithm,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(lifetime) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    appName,
			ID:        userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (string, error) {
	token, err := parseToken(tokenString)
	if err != nil {
		return "", err
	}
	if token.Issuer != appName {
		return "", errors.New("invalid issuer")
	}
	return token.ID, nil
}

func parseToken(tokenString string) (*NexusClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &NexusClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*NexusClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
