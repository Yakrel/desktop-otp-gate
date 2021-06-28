package rand

import (
	"crypto/rand"
	"math/big"
)

const RANDOM_CHARSET = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateRandomString(n int8) (string, error) {
	ret := make([]byte, n)
	var i int8
	for i = 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(RANDOM_CHARSET))))
		if err != nil {
			return "", err
		}
		ret[i] = RANDOM_CHARSET[num.Int64()]
	}
	return string(ret), nil
}
