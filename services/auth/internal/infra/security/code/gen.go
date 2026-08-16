package code

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateNewCode() (string, error) {
	bInt := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, bInt)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
