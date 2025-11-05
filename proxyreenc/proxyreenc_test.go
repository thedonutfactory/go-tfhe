package proxyreenc

import (
	"testing"

	"github.com/thedonutfactory/go-tfhe/key"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
)

func TestPublicKeyEncryption(t *testing.T) {
	secretKey := key.NewSecretKey()
	publicKey := NewPublicKeyLv0(secretKey.KeyLv0)

	// Test encrypting with public key and decrypting with secret key
	messages := []bool{true, false}
	for _, message := range messages {
		ct := publicKey.EncryptBool(message, params.GetTLWELv0().ALPHA)
		decrypted := ct.DecryptBool(secretKey.KeyLv0)

		if decrypted != message {
			t.Errorf("Public key encryption failed: got %v, want %v", decrypted, message)
		}
	}
}

func TestPublicKeyEncryptionMultiple(t *testing.T) {
	secretKey := key.NewSecretKey()
	publicKey := NewPublicKeyLv0(secretKey.KeyLv0)

	correct := 0
	iterations := 100

	for i := 0; i < iterations; i++ {
		message := (i % 2) == 0
		ct := publicKey.EncryptBool(message, params.GetTLWELv0().ALPHA)
		if ct.DecryptBool(secretKey.KeyLv0) == message {
			correct++
		}
	}

	accuracy := float64(correct) / float64(iterations)
	if accuracy < 0.95 {
		t.Errorf("Public key encryption accuracy too low: %.2f%%", accuracy*100)
	}
}

func TestProxyReencryptionAsymmetric(t *testing.T) {
	aliceKey := key.NewSecretKey()
	bobKey := key.NewSecretKey()

	// Bob publishes his public key
	bobPublicKey := NewPublicKeyLv0(bobKey.KeyLv0)

	// Alice generates reencryption key using Bob's PUBLIC key
	reencKey := NewProxyReencryptionKeyAsymmetric(aliceKey.KeyLv0, bobPublicKey)

	// Test both true and false
	messages := []bool{true, false}
	for _, message := range messages {
		aliceCt := tlwe.NewTLWELv0()
		aliceCt.EncryptBool(message, params.GetTLWELv0().ALPHA, aliceKey.KeyLv0)

		// Verify Alice can decrypt
		if aliceCt.DecryptBool(aliceKey.KeyLv0) != message {
			t.Errorf("Alice encryption failed for message %v", message)
		}

		// Reencrypt to Bob's key
		bobCt := ReencryptTLWELv0(aliceCt, reencKey)

		// Verify Bob can decrypt
		decrypted := bobCt.DecryptBool(bobKey.KeyLv0)
		if decrypted != message {
			t.Errorf("Asymmetric proxy reencryption failed: got %v, want %v", decrypted, message)
		}
	}
}

func TestProxyReencryptionSymmetric(t *testing.T) {
	aliceKey := key.NewSecretKey()
	bobKey := key.NewSecretKey()

	// Symmetric mode - requires both secret keys
	reencKey := NewProxyReencryptionKeySymmetric(aliceKey.KeyLv0, bobKey.KeyLv0)

	// Test both true and false
	messages := []bool{true, false}
	for _, message := range messages {
		aliceCt := tlwe.NewTLWELv0()
		aliceCt.EncryptBool(message, params.GetTLWELv0().ALPHA, aliceKey.KeyLv0)

		// Verify Alice can decrypt
		if aliceCt.DecryptBool(aliceKey.KeyLv0) != message {
			t.Errorf("Alice encryption failed for message %v", message)
		}

		// Reencrypt to Bob's key
		bobCt := ReencryptTLWELv0(aliceCt, reencKey)

		// Verify Bob can decrypt
		decrypted := bobCt.DecryptBool(bobKey.KeyLv0)
		if decrypted != message {
			t.Errorf("Symmetric proxy reencryption failed: got %v, want %v", decrypted, message)
		}
	}
}

func TestProxyReencryptionAsymmetricMultiple(t *testing.T) {
	aliceKey := key.NewSecretKey()
	bobKey := key.NewSecretKey()
	bobPublicKey := NewPublicKeyLv0(bobKey.KeyLv0)

	reencKey := NewProxyReencryptionKeyAsymmetric(aliceKey.KeyLv0, bobPublicKey)

	correct := 0
	iterations := 100

	for i := 0; i < iterations; i++ {
		message := (i % 2) == 0

		aliceCt := tlwe.NewTLWELv0()
		aliceCt.EncryptBool(message, params.GetTLWELv0().ALPHA, aliceKey.KeyLv0)

		bobCt := ReencryptTLWELv0(aliceCt, reencKey)

		if bobCt.DecryptBool(bobKey.KeyLv0) == message {
			correct++
		}
	}

	accuracy := float64(correct) / float64(iterations)
	if accuracy < 0.90 {
		t.Errorf("Asymmetric proxy reencryption accuracy too low: %.2f%%", accuracy*100)
	}

	t.Logf("Asymmetric proxy reencryption accuracy: %d/%d (%.1f%%)", correct, iterations, accuracy*100)
}

func TestProxyReencryptionChainAsymmetric(t *testing.T) {
	aliceKey := key.NewSecretKey()
	bobKey := key.NewSecretKey()
	carolKey := key.NewSecretKey()

	bobPublic := NewPublicKeyLv0(bobKey.KeyLv0)
	carolPublic := NewPublicKeyLv0(carolKey.KeyLv0)

	reencKeyAB := NewProxyReencryptionKeyAsymmetric(aliceKey.KeyLv0, bobPublic)
	reencKeyBC := NewProxyReencryptionKeyAsymmetric(bobKey.KeyLv0, carolPublic)

	message := true

	aliceCt := tlwe.NewTLWELv0()
	aliceCt.EncryptBool(message, params.GetTLWELv0().ALPHA, aliceKey.KeyLv0)

	// Alice -> Bob
	bobCt := ReencryptTLWELv0(aliceCt, reencKeyAB)
	if bobCt.DecryptBool(bobKey.KeyLv0) != message {
		t.Errorf("Alice -> Bob reencryption failed")
	}

	// Bob -> Carol
	carolCt := ReencryptTLWELv0(bobCt, reencKeyBC)
	if carolCt.DecryptBool(carolKey.KeyLv0) != message {
		t.Errorf("Bob -> Carol reencryption failed")
	}
}

func TestProxyReencryptionKeyGeneration(t *testing.T) {
	aliceKey := key.NewSecretKey()
	bobKey := key.NewSecretKey()

	reencKey := NewProxyReencryptionKeySymmetric(aliceKey.KeyLv0, bobKey.KeyLv0)

	// Verify the key has the right size
	expectedSize := reencKey.Base * reencKey.T * params.GetTLWELv0().N
	if len(reencKey.KeyEncryptions) != expectedSize {
		t.Errorf("Key size mismatch: got %d, want %d", len(reencKey.KeyEncryptions), expectedSize)
	}

	// Verify structure
	expectedBase := 1 << params.GetTRGSWLv1().BASEBIT
	if reencKey.Base != expectedBase {
		t.Errorf("Base mismatch: got %d, want %d", reencKey.Base, expectedBase)
	}

	if reencKey.T != params.GetTRGSWLv1().IKS_T {
		t.Errorf("T mismatch: got %d, want %d", reencKey.T, params.GetTRGSWLv1().IKS_T)
	}
}

// Benchmark asymmetric key generation
func BenchmarkAsymmetricKeyGeneration(b *testing.B) {
	aliceKey := key.NewSecretKey()
	bobKey := key.NewSecretKey()
	bobPublicKey := NewPublicKeyLv0(bobKey.KeyLv0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewProxyReencryptionKeyAsymmetric(aliceKey.KeyLv0, bobPublicKey)
	}
}

// Benchmark symmetric key generation
func BenchmarkSymmetricKeyGeneration(b *testing.B) {
	aliceKey := key.NewSecretKey()
	bobKey := key.NewSecretKey()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewProxyReencryptionKeySymmetric(aliceKey.KeyLv0, bobKey.KeyLv0)
	}
}

// Benchmark reencryption operation
func BenchmarkReencryption(b *testing.B) {
	aliceKey := key.NewSecretKey()
	bobKey := key.NewSecretKey()

	reencKey := NewProxyReencryptionKeySymmetric(aliceKey.KeyLv0, bobKey.KeyLv0)

	aliceCt := tlwe.NewTLWELv0()
	aliceCt.EncryptBool(true, params.GetTLWELv0().ALPHA, aliceKey.KeyLv0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ReencryptTLWELv0(aliceCt, reencKey)
	}
}

// Benchmark public key generation
func BenchmarkPublicKeyGeneration(b *testing.B) {
	secretKey := key.NewSecretKey()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewPublicKeyLv0(secretKey.KeyLv0)
	}
}

