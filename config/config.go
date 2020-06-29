package config

import (
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
)

type Config struct {
	Host        string
	Port        int

}

func ParseConfig() (*Config, error) {
	host := viper.GetString("Host")
	if host == "" {
		return nil, errors.New("parse host error")
	}

	port := viper.GetInt("Port")
	if port == 0 {
		return nil, errors.New("parse port error")
	}


	return &Config{
		Host:        host,
		Port:        port,
	}, nil
}

func (config *Config) Stringify() ([]byte, error) {
	return json.Marshal(config)
}