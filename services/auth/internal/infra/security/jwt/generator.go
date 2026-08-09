package jwt

import (
	j "MyMessenger/pkg/jwt"
	"MyMessenger/services/auth/internal/config"
	"crypto/rsa"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtGenerator struct {
	prvKey *rsa.PrivateKey
}

func GetNewJwtGenerator(conf *config.AuthConfig) (*jwtGenerator, error) {
	prvKey, err := jwt.ParseRSAPrivateKeyFromPEM(conf.PrvKey)
	if err != nil {
		log.Fatalf("GetNewJwtGenerator(): Trouble with parsing prvKey: %q", err.Error())
		return nil, err
	}
	return &jwtGenerator{prvKey: prvKey}, nil
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
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &claims)
	return token.SignedString(g.prvKey)
}
