package io

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/thedonutfactory/go-tfhe/gates"
)

func WritePubKey(key *gates.PublicKey, filename string) error {
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

func ReadPubKey(filename string) (*gates.PublicKey, error) {
	var key gates.PublicKey
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

func WritePrivKey(key *gates.PrivateKey, filename string) error {
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

func ReadPrivKey(filename string) (*gates.PrivateKey, error) {
	var key gates.PrivateKey
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

func WriteCiphertext(ctxt *gates.Int, filename string) error {
	file, _ := os.Create(filename)
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err := encoder.Encode(ctxt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func ReadCiphertext(filename string) (*gates.Int, error) {
	var ctxt *gates.Int
	dataFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&ctxt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dataFile.Close()
	return ctxt, nil
}
