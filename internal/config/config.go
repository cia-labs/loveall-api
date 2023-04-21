package config

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// Define your configuration fields here
	LogLevel   string `yaml:"logLevel"`
	ServerPort string `yaml:"serverPort"`
}

func Load(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.LogLevel == "" {
		return nil, errors.New("logLevel is required")
	}

	return cfg, nil
}
