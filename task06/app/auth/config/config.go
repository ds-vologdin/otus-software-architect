package config

import (
	"bufio"
	"os"

	"gopkg.in/yaml.v2"
)

// Config of auth service
type Config struct {
	Server         ServerConfig         `yaml:"server"`
	AccountService AccountServiceConfig `yaml:"AccountService"`
	JWT            JWTConfig            `yaml:"jwt"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type JWTConfig struct {
	Algorithm  string `yaml:"algorithm"`
	PrivateKey string `yaml:"privateKey"`
	PublicKey  string `yaml:"publicKey"`
}

type AccountServiceConfig struct {
	URL string `yaml:"url"`
}

// ReadConfig read and decode config file
func ReadConfig(file string) (Config, error) {
	config := Config{}
	f, err := os.Open(file)
	if err != nil {
		return config, err
	}
	r := bufio.NewReader(f)
	d := yaml.NewDecoder(r)
	err = d.Decode(&config)
	return config, err
}
