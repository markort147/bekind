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

var cfg Config

func loadConfig(mockConfig io.Reader) error {
	decoder := yaml.NewDecoder(mockConfig)
	if err := decoder.Decode(&cfg); err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}
	return nil
}

func GetConfig() Config {
	return cfg
}

func FromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	return loadConfig(f)
}
