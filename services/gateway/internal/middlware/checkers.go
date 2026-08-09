package middlware

import (
	"MyMessenger/pkg/jwt"
	"log"
	"strings"
	"time"
)

const (
	tokenPartsCount = 3
)

func (m *Middleware) checkForBanId(claims *jwt.TokenClaims) (time.Time, error) {
	var nilTime time.Time
	id, err := m.checker.GetUserIdFromClaims(claims)
	if err != nil {
		return nilTime, ErrInternal
	}
	expAt, isExists := m.repoBansId.Get(id)
	if !isExists {
		return nilTime, nil
	}
	return expAt, ErrBanned
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

func (m *Middleware) checkToken(token string) error {
	if token == "" {
		log.Printf("checkToken: empty token")
		return ErrBadToken
	}

	if len(token) != m.conf.TokenLenght {
		log.Printf("checkToken: len(token) != m.conf.TokenLenght, %d != %d", len(token), m.conf.TokenLenght)
		return ErrBadToken
	}

	if m.repoTokens.IsExists(token) {
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

	banExpAt, err := m.checkForBanId(claims)
	if err != nil {
		minTime := claims.ExpiresAt.Time
		if banExpAt.Before(minTime) {
			minTime = banExpAt
		}

		m.repoTokens.Put(token, minTime)
		log.Printf("checkToken: new token in repoTokens %q", token)
		return ErrBanned
	}

	return nil
}
