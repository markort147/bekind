package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Log struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"log"`
}

func parseConfig(mockConfig io.Reader) (*Config, error) {
	decoder := yaml.NewDecoder(mockConfig)

	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}
	return &cfg, nil
}

func FromFile(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	return parseConfig(f)
}
