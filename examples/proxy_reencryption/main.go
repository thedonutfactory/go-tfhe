// Proxy Reencryption Example
//
// This example demonstrates how to use LWE proxy reencryption to securely
// delegate access to encrypted data without decryption.
//
// Run with:
//   go run examples/proxy_reencryption/main.go

package main

import (
	"fmt"
	"time"

	"github.com/thedonutfactory/go-tfhe/key"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/proxyreenc"
	"github.com/thedonutfactory/go-tfhe/tlwe"
)

func main() {
	fmt.Println("=== LWE Proxy Reencryption Demo ===\n")

	// Scenario: Alice wants to share encrypted data with Bob
	// without decrypting it, using a semi-trusted proxy

	fmt.Println("1. Setting up keys for Alice and Bob...")
	aliceKey := key.NewSecretKey()
	bobKey := key.NewSecretKey()
	fmt.Println("   ✓ Alice's secret key generated")

	// Bob publishes his public key
	start := time.Now()
	bobPublicKey := proxyreenc.NewPublicKeyLv0(bobKey.KeyLv0)
	pubkeyTime := time.Since(start)
	fmt.Printf("   ✓ Bob's public key generated in %.2fms\n", float64(pubkeyTime.Microseconds())/1000.0)
	fmt.Println("   ✓ Bob shares his public key (safe to publish)\n")

	// Alice encrypts some data
	fmt.Println("2. Alice encrypts her data...")
	messages := []bool{true, false, true, true, false}
	aliceCiphertexts := make([]*tlwe.TLWELv0, len(messages))

	for i, msg := range messages {
		ct := tlwe.NewTLWELv0()
		ct.EncryptBool(msg, params.GetTLWELv0().ALPHA, aliceKey.KeyLv0)
		aliceCiphertexts[i] = ct
	}

	fmt.Println("   Messages encrypted by Alice:")
	for i, msg := range messages {
		fmt.Printf("   - Message %d: %v\n", i+1, msg)
	}
	fmt.Println()

	// Alice generates a proxy reencryption key using Bob's PUBLIC key
	fmt.Println("3. Alice generates a proxy reencryption key (Alice -> Bob)...")
	fmt.Println("   Using ASYMMETRIC mode - Bob's secret key is NOT needed!")
	start = time.Now()
	reencKey := proxyreenc.NewProxyReencryptionKeyAsymmetric(aliceKey.KeyLv0, bobPublicKey)
	keygenTime := time.Since(start)
	fmt.Printf("   ✓ Reencryption key generated in %.2fms\n", float64(keygenTime.Microseconds())/1000.0)
	fmt.Println("   ✓ Alice shares this key with the proxy\n")

	// Proxy reencrypts the data (without learning the plaintext)
	fmt.Println("4. Proxy converts Alice's ciphertexts to Bob's ciphertexts...")
	start = time.Now()
	bobCiphertexts := make([]*tlwe.TLWELv0, len(aliceCiphertexts))
	for i, ct := range aliceCiphertexts {
		bobCiphertexts[i] = proxyreenc.ReencryptTLWELv0(ct, reencKey)
	}
	reencTime := time.Since(start)
	fmt.Printf("   ✓ %d ciphertexts reencrypted in %.2fms\n", len(bobCiphertexts), float64(reencTime.Microseconds())/1000.0)
	fmt.Printf("   ✓ Average time per reencryption: %.2fms\n\n", float64(reencTime.Microseconds())/float64(len(bobCiphertexts))/1000.0)

	// Bob decrypts the reencrypted data
	fmt.Println("5. Bob decrypts the reencrypted data...")
	correct := 0
	decryptedMessages := make([]bool, len(bobCiphertexts))

	for i, ct := range bobCiphertexts {
		decryptedMessages[i] = ct.DecryptBool(bobKey.KeyLv0)
	}

	fmt.Println("   Decrypted messages:")
	for i, original := range messages {
		decrypted := decryptedMessages[i]
		status := "✗"
		if original == decrypted {
			correct++
			status = "✓"
		}
		fmt.Printf("   %s Message %d: %v (original: %v)\n", status, i+1, decrypted, original)
	}
	fmt.Println()

	fmt.Println("=== Results ===")
	accuracy := float64(correct) / float64(len(messages)) * 100.0
	fmt.Printf("Accuracy: %d/%d (%.1f%%)\n", correct, len(messages), accuracy)
	fmt.Println()

	// Demonstrate multi-hop reencryption: Alice -> Bob -> Carol
	fmt.Println("\n=== Multi-Hop Reencryption Demo (Asymmetric) ===\n")
	fmt.Println("Demonstrating a chain: Alice -> Bob -> Carol")
	fmt.Println("Each party only needs the next party's PUBLIC key\n")

	carolKey := key.NewSecretKey()
	carolPublicKey := proxyreenc.NewPublicKeyLv0(carolKey.KeyLv0)
	fmt.Println("1. Carol's keys generated and public key published")

	reencKeyBC := proxyreenc.NewProxyReencryptionKeyAsymmetric(bobKey.KeyLv0, carolPublicKey)
	fmt.Println("2. Generated reencryption key (Bob -> Carol) using Carol's PUBLIC key")

	testMessage := true
	aliceCt := tlwe.NewTLWELv0()
	aliceCt.EncryptBool(testMessage, params.GetTLWELv0().ALPHA, aliceKey.KeyLv0)
	fmt.Printf("3. Alice encrypts message: %v\n", testMessage)

	bobCt := proxyreenc.ReencryptTLWELv0(aliceCt, reencKey)
	fmt.Println("4. Proxy reencrypts Alice -> Bob")
	bobDecrypted := bobCt.DecryptBool(bobKey.KeyLv0)
	bobStatus := "✗"
	if bobDecrypted == testMessage {
		bobStatus = "✓"
	}
	fmt.Printf("   Bob decrypts: %v %s\n", bobDecrypted, bobStatus)

	carolCt := proxyreenc.ReencryptTLWELv0(bobCt, reencKeyBC)
	fmt.Println("5. Proxy reencrypts Bob -> Carol")
	carolDecrypted := carolCt.DecryptBool(carolKey.KeyLv0)
	carolStatus := "✗"
	if carolDecrypted == testMessage {
		carolStatus = "✓"
	}
	fmt.Printf("   Carol decrypts: %v %s\n", carolDecrypted, carolStatus)

	fmt.Println()
	fmt.Println("=== Security Notes ===")
	fmt.Println("• The proxy never learns the plaintext")
	fmt.Println("• Bob's secret key is NEVER shared - only his public key is used")
	fmt.Println("• The reencryption key only works in one direction")
	fmt.Println("• Each reencryption adds a small amount of noise")
	fmt.Println("• The scheme is unidirectional (Alice->Bob key ≠ Bob->Alice key)")
	fmt.Println("• True asymmetric proxy reencryption with LWE-based public keys")

	fmt.Println("\n=== Performance Summary ===")
	fmt.Printf("Bob's public key generation: %.2fms\n", float64(pubkeyTime.Microseconds())/1000.0)
	fmt.Printf("Reencryption key generation: %.2fms\n", float64(keygenTime.Microseconds())/1000.0)
	fmt.Printf("Average reencryption time: %.2fms\n", float64(reencTime.Microseconds())/float64(len(bobCiphertexts))/1000.0)
}

