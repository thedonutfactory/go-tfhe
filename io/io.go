package io

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/thedonutfactory/go-tfhe/core"
)

func WritePubKey(key *core.PublicKey, filename string) error {
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

func ReadPubKey(filename string) (*core.PublicKey, error) {
	var key core.PublicKey
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

func WritePrivKey(key *core.PrivateKey, filename string) error {
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

func ReadPrivKey(filename string) (*core.PrivateKey, error) {
	var key core.PrivateKey
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
