package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Host string `json:"Host"`
	MaxHttpConnections  int `json:"MaxHttpConnections"`
	Target  int `json:"Target"`
	Timeout  int `json:"Timeout"` // in seconds
	PayloadsPerConnection  int `json:"PayloadsPerConnection"`
}

func (c *Config) GetTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

// client is responsible for duplicate check :)

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
	return &config, nil
}
