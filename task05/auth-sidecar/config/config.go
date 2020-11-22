package config

import (
	"bufio"
	"os"

	"gopkg.in/yaml.v2"
)

// Config of auth service
type Config struct {
	Server      ServerConfig `yaml:"server"`
	JWT         JWTConfig    `yaml:"jwt"`
	Target      TargetConfig `yaml:"target"`
	ExcludeAuth []Request    `yaml:"excludeAuth"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type JWTConfig struct {
	Algorithm string `yaml:"algorithm"`
	PublicKey string `yaml:"publicKey"`
}

type TargetConfig struct {
	URL string `yaml:"url"`
}

type Request struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
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
