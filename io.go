package tfhe

import (
	"encoding/gob"
	"fmt"
	"os"
)

func WritePubKey(key *PublicKey, filename string) error {
	file, _ := os.Create(filename)
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err := encoder.Encode(key)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func ReadPubKey(filename string) (*PublicKey, error) {
	var key PublicKey
	dataFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dataFile.Close()
	return &key, nil
}

func WritePrivKey(key *PrivateKey, filename string) error {
	file, _ := os.Create(filename)
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err := encoder.Encode(key)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func ReadPrivKey(filename string) (*PrivateKey, error) {
	var key PrivateKey
	dataFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dataFile.Close()
	return &key, nil
}
