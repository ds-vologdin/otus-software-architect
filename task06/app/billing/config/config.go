package config

import (
	"bufio"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

// ReadConfig read a config from a file
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
