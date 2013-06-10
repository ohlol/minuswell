package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Transports map[string]map[string]interface{}
	Files      map[string]map[string]interface{}
}

func ReadConfig(filename string) (*Config, error) {
	var config Config

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	dec := json.NewDecoder(file)

	err = dec.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
