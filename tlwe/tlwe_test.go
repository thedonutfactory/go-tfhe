package tlwe_test

import (
	"testing"

	"github.com/thedonutfactory/go-tfhe/key"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
)

func TestTLWELv0EncryptDecrypt(t *testing.T) {
	sk := key.NewSecretKey()

	testCases := []bool{true, false}

	for _, val := range testCases {
		ct := tlwe.NewTLWELv0().EncryptBool(val, params.GetTLWELv0().ALPHA, sk.KeyLv0)
		dec := ct.DecryptBool(sk.KeyLv0)

		if dec != val {
			t.Errorf("Encrypt/Decrypt(%v) = %v", val, dec)
		}
	}
}

func TestTLWELv0EncryptDecryptMultiple(t *testing.T) {
	sk := key.NewSecretKey()
	trials := 100
	correct := 0

	for i := 0; i < trials; i++ {
		val := i%2 == 0
		ct := tlwe.NewTLWELv0().EncryptBool(val, params.GetTLWELv0().ALPHA, sk.KeyLv0)
		dec := ct.DecryptBool(sk.KeyLv0)

		if dec == val {
			correct++
		}
	}

	if correct != trials {
		t.Errorf("Correctness: %d/%d (%.1f%%)", correct, trials, float64(correct)/float64(trials)*100)
	}
}

func TestTLWELv0Add(t *testing.T) {
	sk := key.NewSecretKey()

	ct1 := tlwe.NewTLWELv0().EncryptBool(true, params.GetTLWELv0().ALPHA, sk.KeyLv0)
	ct2 := tlwe.NewTLWELv0().EncryptBool(false, params.GetTLWELv0().ALPHA, sk.KeyLv0)

	sum := ct1.Add(ct2)

	// Addition of true + false should still be decryptable
	// (though the semantic meaning depends on the circuit)
	_ = sum.DecryptBool(sk.KeyLv0)
}

func TestTLWELv0Neg(t *testing.T) {
	sk := key.NewSecretKey()

	ct := tlwe.NewTLWELv0().EncryptBool(true, params.GetTLWELv0().ALPHA, sk.KeyLv0)
	negCt := ct.Neg()
	dec := negCt.DecryptBool(sk.KeyLv0)

	// Negation of true (0.125) should give false (-0.125)
	if dec != false {
		t.Errorf("Neg(true) = %v, expected false", dec)
	}
}

func TestTLWELv1EncryptDecrypt(t *testing.T) {
	sk := key.NewSecretKey()

	testCases := []bool{true, false}

	for _, val := range testCases {
		ct := tlwe.NewTLWELv1().EncryptBool(val, params.GetTLWELv1().ALPHA, sk.KeyLv1)
		dec := ct.DecryptBool(sk.KeyLv1)

		if dec != val {
			t.Errorf("TLWELv1 Encrypt/Decrypt(%v) = %v", val, dec)
		}
	}
}
