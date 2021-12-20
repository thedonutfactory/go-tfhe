package io

import (
	"os"
	"testing"

	"github.com/thedonutfactory/go-tfhe/core"
)

func TestWriteKeys(t *testing.T) {
	// generate params
	params := core.NewDefaultGateBootstrappingParameters(100)

	pubKey, privKey := core.GenerateKeys(params)
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
