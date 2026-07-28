package service

import (
	j "MyMessenger/pkg/jwt"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenToken(userId uuid.UUID, tTl time.Duration, prvKey *rsa.PrivateKey, tType j.TokenType) (string, error) {
	claims := j.TokenClaims{
		Ttype: tType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userId.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(tTl)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &claims)
	return token.SignedString(prvKey)
}
