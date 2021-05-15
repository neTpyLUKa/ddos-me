package client_config

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
	Address string `json:"Address"`
	PublicKey string `json:"PublicKey"` // contract public key
}

func (c *Config) GetTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

// todo try some http libraries for better performance? https://github.com/gojek/heimdall

// client is responsible for duplicate check :)
// https://github.com/lucas-clemente/quic-go/blob/master/example/main.go
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
