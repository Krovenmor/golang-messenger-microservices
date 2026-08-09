package argon2id

import (
	"MyMessenger/services/auth/internal/config"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	hashPattern = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	hashParts   = 6
)

var (
	ErrRandTrouble     = errors.New("rand trouble")
	ErrInvalidHash     = errors.New("invalid hash")
	ErrInvalidVersion  = errors.New("incompatible argon2 version")
	ErrInvalidParams   = errors.New("invalid hash parameters")
	ErrDecodingTrouble = errors.New("decoding trouble")

	ErrComparingFailed = errors.New("false comparing")
)

type argon2idHasher struct {
	SaltLength  int
	Iterations  uint32
	Memory      uint32
	Parallelism uint8
	KeyLength   uint32
}

func NewArgon2idHasher(conf *config.HashConfig) *argon2idHasher {
	return &argon2idHasher{
		SaltLength:  conf.SaltLength,
		Iterations:  uint32(conf.Iterations),
		Memory:      uint32(conf.Memory),
		Parallelism: uint8(conf.Parallelism),
		KeyLength:   uint32(conf.KeyLength),
	}
}

func (h *argon2idHasher) Hash(password string) (string, error) {
	salt := make([]byte, h.SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		log.Printf("Hash: trouble with rand.Read(salt), err: %q", err)
		return "", ErrRandTrouble
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.Iterations,
		h.Memory,
		h.Parallelism,
		h.KeyLength,
	)

	return fmt.Sprintf(
		hashPattern,
		argon2.Version,
		h.Memory,
		h.Iterations,
		h.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func (h *argon2idHasher) Comp(password, hash string) error {
	parts := strings.Split(hash, "$")
	if len(parts) != hashParts {
		log.Printf("Comp: invalid hash, hash: %q", hash)
		return ErrInvalidHash
	}
	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		log.Printf("Comp: invalid version: %v, we have: %v", version, argon2.Version)
		return ErrInvalidVersion
	}

	var Iterations, Memory uint32
	var Parallelism uint8

	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &Memory, &Iterations, &Parallelism)
	if err != nil {
		log.Printf("Comp: invalid params, err: %q", err)
		return ErrInvalidParams
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		log.Printf("Comp: trouble with decoding salt, err: %q", err)
		return ErrDecodingTrouble
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		log.Printf("Comp: trouble with decoding hash, err: %q", err)
		return ErrDecodingTrouble
	}

	KeyLength := uint32(len(decodedHash))

	comparisonHash := argon2.IDKey(
		[]byte(password),
		salt,
		Iterations,
		Memory,
		Parallelism,
		KeyLength,
	)

	if subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 0 {
		return ErrComparingFailed
	}

	return nil
}
