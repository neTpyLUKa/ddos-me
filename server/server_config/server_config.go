package server_config

import (
	"encoding/json"
	"github.com/caarlos0/env"
	"log"
	"os"
)

type Config struct {
	MaxHttpConnections  int `json:"MaxHttpConnections"`
	PrivateKey string `json:"PrivateKey"`
}

// results file ?
// retry count ?
// http_timeout ?
// public server key ?

func ParseConfig(configFile string) (*Config, error) {
	configReader, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	jsonParser := json.NewDecoder(configReader)
	config := Config{}
	if err := jsonParser.Decode(&config); err != nil {
		return nil, err
	}
	if err := env.Parse(&config); err != nil {
		log.Fatalln("Error processing env variables", err.Error())
	}
	return &config, nil
}
