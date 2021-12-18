package tfhe

import (
	"os"
	"testing"
)

func TestWriteKeys(t *testing.T) {
	// generate params
	params := NewDefaultGateBootstrappingParameters(100)

	pubKey, privKey := GenerateKeys(params)
	defer func() {
		os.Remove("private.key")
		os.Remove("public.key")
	}()

	err := WritePrivKey(privKey, "private.key")
	if err != nil {
		t.Errorf("Could not serialize the private key to file")
	}
	err = WritePubKey(pubKey, "public.key")
	if err != nil {
		t.Errorf("Could not serialize the public key to file")
	}
}
