package tfhe

import (
	"encoding/gob"
	"fmt"
	"os"
)

func SavePubKey(key *PubKey, filename string) error {
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

func OpenPubKey(filename string) (*PubKey, error) {
	var key PubKey
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

func SavePrivKey(key *PriKey, filename string) error {
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

func OpenPrivKey(filename string) (*PriKey, error) {
	var key PriKey
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
