package middlware

import (
	"MyMessenger/pkg/jwt"
	"MyMessenger/pkg/utils"
	"context"
	"log"
	"strings"
	"time"
)

const (
	tokenPartsCount = 3
)

func (m *Middleware) checkForBanId(ctx context.Context, claims *jwt.TokenClaims) (time.Duration, error) {
	id, err := m.checker.GetUserIdFromClaims(claims)
	if err != nil {
		return -1, ErrInternal
	}
	ttl, isExists := m.repoBansId.Get(ctx, id.String())
	if !isExists {
		return -1, nil
	}
	return ttl, ErrBanned
}

func (m *Middleware) checkTokenParts(token string) error {
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) != tokenPartsCount {
		log.Printf("checkTokenParts: not valid access token, Parts len %d != %d", len(tokenParts), tokenPartsCount)
		return ErrBadToken
	}

	if len(tokenParts[0]) != m.conf.HeaderLength {
		log.Printf("checkTokenParts: not valid access token, Header len %d != %d", len(tokenParts[0]), m.conf.HeaderLength)
		return ErrBadToken
	}

	if len(tokenParts[1]) != m.conf.PayloadLength {
		log.Printf("checkTokenParts: not valid access token, Payload len %d != %d", len(tokenParts[1]), m.conf.PayloadLength)
		return ErrBadToken
	}

	if len(tokenParts[2]) != m.conf.SignatureLength {
		log.Printf("checkTokenParts: not valid access token, Signature len %d != %d", len(tokenParts[2]), m.conf.SignatureLength)
		return ErrBadToken
	}
	return nil
}

func (m *Middleware) checkToken(ctx context.Context, token string) error {
	if token == "" {
		log.Printf("checkToken: empty token")
		return ErrBadToken
	}

	if len(token) != m.conf.TokenLenght {
		log.Printf("checkToken: len(token) != m.conf.TokenLenght, %d != %d", len(token), m.conf.TokenLenght)
		return ErrBadToken
	}

	if m.repoTokens.IsExists(ctx, token) {
		log.Printf("checkToken: banned token %q in repoTokens", token)
		return ErrBanned
	}

	err := m.checkTokenParts(token)
	if err != nil {
		return err
	}

	claims, err := m.checker.IsValidAccess(token)
	if err != nil {
		log.Printf("checkToken: not valid access token, err: %q", err)
		return ErrBadToken
	}

	banTtl, err := m.checkForBanId(ctx, claims)
	if err != nil {
		hashed := utils.Hash(token)
		err := m.repoTokens.Put(ctx, hashed, banTtl)
		if err != nil {
			log.Printf("checkToken: trouble with saving token, token: %q, err: %q", hashed, err)
		} else {
			log.Printf("checkToken: new token in repoTokens %q", hashed)
		}
		return ErrBanned
	}

	return nil
}
