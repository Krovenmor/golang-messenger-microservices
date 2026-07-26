package service

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenToken(userId uuid.UUID, tTl time.Duration, prvKey *rsa.PrivateKey) (string, error) {
	mc := jwt.RegisteredClaims{
		ID:        uuid.New().String(),
		Subject:   userId.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(tTl)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, mc)
	return token.SignedString(prvKey)
}

func IsValidToken(cToken string, pubKey *rsa.PublicKey) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(cToken, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token is not valid")
	}

	return claims, nil
}
