package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Outputs map[string]map[string]interface{}
	Files   map[string]FilesConfig
}

type FilesConfig struct {
	Type   string
	Tags   []string
	Fields map[string]interface{}
}

func ReadConfig(filename string) (*Config, error) {
	var config Config

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
