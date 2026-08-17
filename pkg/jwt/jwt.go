package jwt

import (
	"MyMessenger/pkg/config"
	"crypto/ed25519"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
)

type TokenClaims struct {
	Ttype TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type JWTChecker struct {
	pubKey    ed25519.PublicKey
	aTokenLen int
}

func NewJwtCheckerConf(conf config.JwtCheckerConf) (*JWTChecker, error) {
	pubKey, err := jwt.ParseEdPublicKeyFromPEM(conf.PubKey)
	if err != nil {
		log.Printf("NewJwtCheckerConf(): Trouble with parsing pubKey: %q", err.Error())
		return nil, err
	}
	edKey, ok := pubKey.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("GetNewJwtGenerator: key is not of type ed25519.PrivateKey")
	}
	return &JWTChecker{pubKey: edKey, aTokenLen: conf.ATokenLen}, nil
}

func (j *JWTChecker) IsValidToken(cToken string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(cToken, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.pubKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}

func (j *JWTChecker) IsValidAccess(aToken string) (*TokenClaims, error) {
	if len(aToken) != j.aTokenLen {
		return nil, fmt.Errorf("wrong len, %d != %d", len(aToken), j.aTokenLen)
	}
	claims, err := j.IsValidToken(aToken)
	if err != nil {
		return nil, err
	}
	if claims.Ttype != Access {
		return nil, fmt.Errorf("not an Access Token")
	}
	return claims, nil
}

func (j *JWTChecker) GetUserIdFromClaims(claims *TokenClaims) (uuid.UUID, error) {
	if claims == nil {
		return uuid.Nil, fmt.Errorf("claims == nil")
	}
	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
