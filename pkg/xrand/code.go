package xrand

import (
	"crypto/rand"
	"math"
	"math/big"
)

func GenerateRandomCode(maxDigits int) (string, error) {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))),
	)
	if err != nil {
		return "", err
	}
	return bi.String(), nil
}
