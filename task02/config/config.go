package config

import (
	"bufio"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
}

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
