package tfhe

import (
	"fmt"
	"os"
	"testing"
)

func TestWritePublicKey(t *testing.T) {
	// generate params
	var minimumLambda int32 = 100
	params := NewDefaultGateBootstrappingParameters(minimumLambda)

	pubKey, privKey := GenerateKeys(params)
	defer func() {
		os.Remove("private.key")
		os.Remove("public.key")
	}()

	fmt.Printf("%+v\n", pubKey.Bk.Bk)
	fmt.Printf("%+v\n", pubKey.Bk.AccumParams)
	fmt.Printf("%+v\n", pubKey.Bk.Ks)

	err := WritePrivKey(privKey, "private.key")
	if err != nil {
		t.Errorf("Could not serialize the private key to file")
	}
	err = WritePubKey(pubKey, "public.key")
	if err != nil {
		t.Errorf("Could not serialize the public key to file")
	}
}
