// Package proxyreenc implements LWE-based proxy reencryption for TFHE.
//
// This package provides proxy reencryption functionality that allows secure
// transformation of ciphertexts from one secret key to another without decryption.
//
// Proxy reencryption enables a semi-trusted proxy to convert a ciphertext encrypted
// under one key (delegator) to a ciphertext encrypted under another key (delegatee)
// without learning the plaintext.
//
// # Two Modes of Operation
//
// Asymmetric Mode (Recommended): Alice generates a reencryption key using her
// secret key and Bob's public key. Bob never shares his secret key.
//
// Symmetric Mode (Trusted Scenarios): When both secret keys are available,
// such as for single-party key rotation.
//
// # Example (Asymmetric Mode)
//
//	// Bob publishes his public key
//	bobPublicKey := proxyreenc.NewPublicKeyLv0(bobSecretKey.KeyLv0)
//
//	// Alice generates reencryption key using Bob's PUBLIC key
//	reencKey := proxyreenc.NewProxyReencryptionKeyAsymmetric(aliceSecretKey.KeyLv0, bobPublicKey)
//
//	// Proxy transforms ciphertext (learns nothing about plaintext)
//	bobCt := proxyreenc.ReencryptTLWELv0(aliceCt, reencKey)
//
//	// Bob decrypts with his secret key
//	plaintext := bobCt.DecryptBool(bobSecretKey.KeyLv0)
//
// # Security
//
// - The proxy learns nothing about the plaintext during reencryption
// - Reencryption keys are unidirectional (Alice→Bob ≠ Bob→Alice)
// - In asymmetric mode, Bob's secret key is never exposed
// - Based on the hardness of the Learning With Errors (LWE) problem
// - 128-bit post-quantum security
package proxyreenc

import (
	"math/rand"

	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
	"github.com/thedonutfactory/go-tfhe/utils"
)

// PublicKeyLv0 represents an LWE public key for asymmetric encryption.
//
// A public key consists of encryptions of zero under the secret key.
// This allows anyone to encrypt messages without knowing the secret key.
//
// The public key can be safely shared without revealing the secret key
// due to the semantic security of LWE.
type PublicKeyLv0 struct {
	Encryptions []*tlwe.TLWELv0 // Encryptions of zero
}

// NewPublicKeyLv0 generates a new public key from a secret key.
//
// Creates encryptions of zero that can be used for encryption without
// revealing the secret key.
//
// Parameters:
//   - secretKey: The secret key to generate a public key for
//
// Returns a public key that can be safely shared.
func NewPublicKeyLv0(secretKey []params.Torus) *PublicKeyLv0 {
	return NewPublicKeyLv0WithParams(secretKey, params.GetTLWELv0().N*2, params.GetTLWELv0().ALPHA)
}

// NewPublicKeyLv0WithParams generates a public key with custom parameters.
//
// Parameters:
//   - secretKey: The secret key
//   - size: Number of zero encryptions to generate (larger = more security)
//   - alpha: Noise parameter for encryptions
func NewPublicKeyLv0WithParams(secretKey []params.Torus, size int, alpha float64) *PublicKeyLv0 {
	encryptions := make([]*tlwe.TLWELv0, size)

	// Generate encryptions of zero
	for i := 0; i < size; i++ {
		ct := tlwe.NewTLWELv0()
		ct.EncryptF64(0.0, alpha, secretKey)
		encryptions[i] = ct
	}

	return &PublicKeyLv0{
		Encryptions: encryptions,
	}
}

// EncryptF64 encrypts a value using the public key.
//
// This allows encryption without the secret key by combining
// the pre-computed zero encryptions.
//
// Parameters:
//   - plaintext: Value to encrypt (as float64)
//   - alpha: Additional noise parameter
//
// Returns a TLWELv0 ciphertext encrypting the plaintext.
func (pk *PublicKeyLv0) EncryptF64(plaintext float64, alpha float64) *tlwe.TLWELv0 {
	rng := rand.New(rand.NewSource(rand.Int63()))
	result := tlwe.NewTLWELv0()

	// Add the plaintext to b
	plaintextTorus := utils.F64ToTorus(plaintext)
	result.SetB(plaintextTorus)

	n := params.GetTLWELv0().N

	// Randomly combine encryptions of zero
	for _, enc := range pk.Encryptions {
		if rng.Intn(2) == 1 {
			// Add or subtract randomly
			if rng.Intn(2) == 1 {
				for i := 0; i <= n; i++ {
					result.P[i] = result.P[i] + enc.P[i]
				}
			} else {
				for i := 0; i <= n; i++ {
					result.P[i] = result.P[i] - enc.P[i]
				}
			}
		}
	}

	// Add fresh noise
	noise := utils.GaussianF64(0.0, alpha, rng)
	result.SetB(result.B() + noise)

	return result
}

// EncryptBool encrypts a boolean using the public key.
//
// Parameters:
//   - plaintext: Boolean value to encrypt
//   - alpha: Noise parameter
//
// Returns a TLWELv0 ciphertext encrypting the boolean.
func (pk *PublicKeyLv0) EncryptBool(plaintext bool, alpha float64) *tlwe.TLWELv0 {
	var p float64
	if plaintext {
		p = 0.125
	} else {
		p = -0.125
	}
	return pk.EncryptF64(p, alpha)
}

// ProxyReencryptionKey stores the reencryption key from one secret key to another.
//
// This key allows converting ciphertexts encrypted under keyFrom to
// ciphertexts encrypted under keyTo. It uses a decomposition-based
// approach similar to key switching in TFHE.
type ProxyReencryptionKey struct {
	KeyEncryptions []*tlwe.TLWELv0 // Decomposed encryptions for key switching
	Base           int             // Base for decomposition (typically 1 << BASEBIT)
	T              int             // Number of decomposition levels
}

// NewProxyReencryptionKeyAsymmetric generates a proxy reencryption key using asymmetric mode (RECOMMENDED).
//
// Alice generates a reencryption key using her secret key and Bob's public key.
// Bob never needs to share his secret key with Alice.
//
// Parameters:
//   - keyFrom: Alice's secret key (delegator)
//   - publicKeyTo: Bob's public key (delegatee)
//
// Returns a proxy reencryption key from Alice to Bob.
//
// # Security
//
// This is the secure way to generate a reencryption key. Bob's secret key
// is never exposed, only his public key is needed.
func NewProxyReencryptionKeyAsymmetric(keyFrom []params.Torus, publicKeyTo *PublicKeyLv0) *ProxyReencryptionKey {
	return NewProxyReencryptionKeyAsymmetricWithParams(
		keyFrom,
		publicKeyTo,
		params.KSKAlpha(),
		params.GetTRGSWLv1().BASEBIT,
		params.GetTRGSWLv1().IKS_T,
	)
}

// NewProxyReencryptionKeyAsymmetricWithParams generates a proxy reencryption key with custom parameters.
func NewProxyReencryptionKeyAsymmetricWithParams(
	keyFrom []params.Torus,
	publicKeyTo *PublicKeyLv0,
	alpha float64,
	basebit int,
	t int,
) *ProxyReencryptionKey {
	base := 1 << basebit
	n := params.GetTLWELv0().N

	keyEncryptions := make([]*tlwe.TLWELv0, base*t*n)

	// Initialize all to zero
	for i := range keyEncryptions {
		keyEncryptions[i] = tlwe.NewTLWELv0()
	}

	// Generate decomposed encryptions using the PUBLIC key
	for i := 0; i < n; i++ {
		for j := 0; j < t; j++ {
			for k := 0; k < base; k++ {
				if k == 0 {
					continue // Skip k=0 as it contributes nothing
				}

				// Encrypt k * keyFrom[i] / 2^((j+1)*basebit) using Bob's PUBLIC key
				shiftAmount := (j + 1) * basebit
				p := (float64(k) * float64(keyFrom[i])) / float64(uint32(1)<<shiftAmount)
				idx := (base * t * i) + (base * j) + k

				// Use public key encryption instead of secret key
				keyEncryptions[idx] = publicKeyTo.EncryptF64(p, alpha)
			}
		}
	}

	return &ProxyReencryptionKey{
		KeyEncryptions: keyEncryptions,
		Base:           base,
		T:              t,
	}
}

// NewProxyReencryptionKeySymmetric generates a proxy reencryption key using symmetric mode.
//
// This requires both secret keys and should only be used in trusted scenarios
// like single-party key rotation or when both parties trust each other.
//
// Parameters:
//   - keyFrom: Source secret key
//   - keyTo: Target secret key
//
// Returns a proxy reencryption key from source to target.
//
// # Security Warning
//
// This mode requires access to both secret keys. For true delegation where
// Bob doesn't share his secret key, use NewProxyReencryptionKeyAsymmetric instead.
func NewProxyReencryptionKeySymmetric(keyFrom []params.Torus, keyTo []params.Torus) *ProxyReencryptionKey {
	return NewProxyReencryptionKeySymmetricWithParams(
		keyFrom,
		keyTo,
		params.KSKAlpha(),
		params.GetTRGSWLv1().BASEBIT,
		params.GetTRGSWLv1().IKS_T,
	)
}

// NewProxyReencryptionKeySymmetricWithParams generates a symmetric mode key with custom parameters.
func NewProxyReencryptionKeySymmetricWithParams(
	keyFrom []params.Torus,
	keyTo []params.Torus,
	alpha float64,
	basebit int,
	t int,
) *ProxyReencryptionKey {
	base := 1 << basebit
	n := params.GetTLWELv0().N

	keyEncryptions := make([]*tlwe.TLWELv0, base*t*n)

	// Initialize all to zero
	for i := range keyEncryptions {
		keyEncryptions[i] = tlwe.NewTLWELv0()
	}

	// Generate decomposed encryptions similar to key switching key
	for i := 0; i < n; i++ {
		for j := 0; j < t; j++ {
			for k := 0; k < base; k++ {
				if k == 0 {
					continue // Skip k=0 as it contributes nothing
				}

				// Encrypt k * keyFrom[i] / 2^((j+1)*basebit)
				shiftAmount := (j + 1) * basebit
				p := (float64(k) * float64(keyFrom[i])) / float64(uint32(1)<<shiftAmount)
				idx := (base * t * i) + (base * j) + k

				keyEncryptions[idx].EncryptF64(p, alpha, keyTo)
			}
		}
	}

	return &ProxyReencryptionKey{
		KeyEncryptions: keyEncryptions,
		Base:           base,
		T:              t,
	}
}

// ReencryptTLWELv0 reencrypts a TLWELv0 ciphertext from one key to another.
//
// Converts a ciphertext encrypted under the source key (embedded in the
// reencryption key) to a ciphertext encrypted under the target key.
//
// Parameters:
//   - ctFrom: Ciphertext encrypted under the source key
//   - reencKey: Proxy reencryption key from source to target
//
// Returns a new ciphertext encrypting the same plaintext under the target key.
//
// # Algorithm
//
// Uses a decomposition-based approach similar to identity key switching:
// 1. Start with the b value from the source ciphertext
// 2. For each coefficient a[i] in the source:
//   - Decompose a[i] into digits in base reencKey.Base
//   - Subtract the corresponding pre-computed encrypted values
// 3. Result is an encryption of the same message under the target key
func ReencryptTLWELv0(ctFrom *tlwe.TLWELv0, reencKey *ProxyReencryptionKey) *tlwe.TLWELv0 {
	n := params.GetTLWELv0().N

	// Calculate basebit from base
	basebit := 0
	for (1 << basebit) < reencKey.Base {
		basebit++
	}

	base := reencKey.Base
	t := reencKey.T

	result := tlwe.NewTLWELv0()

	// Start with the b value from the source ciphertext
	result.SetB(ctFrom.B())

	// Precision offset for rounding (similar to identity key switching)
	precOffset := params.Torus(1) << (32 - (1 + basebit*t))

	// Process each coefficient of the source ciphertext
	for i := 0; i < n; i++ {
		// Add precision offset for rounding
		aBar := ctFrom.P[i] + precOffset

		// Decompose into t levels
		for j := 0; j < t; j++ {
			// Extract the j-th digit in base `base`
			shift := 32 - (j+1)*basebit
			mask := (params.Torus(1) << basebit) - 1
			k := (aBar >> shift) & mask

			if k != 0 {
				// Index into the reencryption key
				idx := (base * t * i) + (base * j) + int(k)

				// Subtract the pre-computed encryption
				for x := 0; x <= n; x++ {
					result.P[x] = result.P[x] - reencKey.KeyEncryptions[idx].P[x]
				}
			}
		}
	}

	return result
}

