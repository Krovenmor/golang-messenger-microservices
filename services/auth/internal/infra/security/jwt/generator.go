package jwt

import (
	j "MyMessenger/pkg/jwt"
	"MyMessenger/services/auth/internal/config"
	"crypto/ed25519"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtGenerator struct {
	prvKey ed25519.PrivateKey
}

func GetNewJwtGenerator(conf *config.AuthConfig) (*jwtGenerator, error) {
	prvKey, err := jwt.ParseEdPrivateKeyFromPEM(conf.PrvKey)
	if err != nil {
		log.Printf("GetNewJwtGenerator(): Trouble with parsing prvKey: %q", err.Error())
		return nil, err
	}
	edKey, ok := prvKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("GetNewJwtGenerator: key is not of type ed25519.PrivateKey")
	}
	return &jwtGenerator{prvKey: edKey}, nil
}

func (g *jwtGenerator) GenToken(userId uuid.UUID, tTl time.Duration, tType j.TokenType) (string, error) {
	claims := j.TokenClaims{
		Ttype: tType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userId.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(tTl)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, &claims)
	return token.SignedString(g.prvKey)
}
