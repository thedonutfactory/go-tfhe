package key

import (
	"math/rand"

	"github.com/thedonutfactory/go-tfhe/params"
)

// SecretKey contains the secret keys for both levels
type SecretKey struct {
	KeyLv0 []params.Torus
	KeyLv1 []params.Torus
}

// NewSecretKey generates a new secret key
func NewSecretKey() *SecretKey {
	rng := rand.New(rand.NewSource(rand.Int63()))

	lv0N := params.GetTLWELv0().N
	lv1N := params.GetTLWELv1().N

	keyLv0 := make([]params.Torus, lv0N)
	keyLv1 := make([]params.Torus, lv1N)

	for i := 0; i < lv0N; i++ {
		if rng.Intn(2) == 1 {
			keyLv0[i] = 1
		} else {
			keyLv0[i] = 0
		}
	}

	for i := 0; i < lv1N; i++ {
		if rng.Intn(2) == 1 {
			keyLv1[i] = 1
		} else {
			keyLv1[i] = 0
		}
	}

	return &SecretKey{
		KeyLv0: keyLv0,
		KeyLv1: keyLv1,
	}
}
