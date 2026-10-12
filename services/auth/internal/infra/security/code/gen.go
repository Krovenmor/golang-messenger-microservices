package code

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

const codeLength = 6

var codeUpperBound = big.NewInt(int64(math.Pow10(codeLength)))

func GenerateNewCode() (string, error) {
	n, err := rand.Int(rand.Reader, codeUpperBound)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", codeLength, n.Int64()), nil
}
