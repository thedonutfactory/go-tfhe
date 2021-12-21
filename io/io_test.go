package io

import (
	"os"
	"testing"

	"github.com/thedonutfactory/go-tfhe/gates"
)

func TestWriteKeys(t *testing.T) {
	// generate params
	params := gates.DefaultGateBootstrappingParameters(100)

	pubKey, privKey := params.GenerateKeys()
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
