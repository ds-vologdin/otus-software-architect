package config

import (
	"bufio"
	"os"

	"gopkg.in/yaml.v2"
)

// Config of auth service
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	AccountService AccountServiceConfig `yaml:"AccountService"`
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
