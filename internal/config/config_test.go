package config

import (
	"io"
	"os"
	"strings"
	"testing"
)

type mockReader struct {
	data string
	read bool
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	if m.read {
		return 0, io.EOF
	}
	n = copy(p, m.data)
	m.read = true
	return n, nil
}

func TestLoadConfig(t *testing.T) {

	t.Run("testing private func loadConfig", func(t *testing.T) {
		mockConfig := &mockReader{
			data: `server:
  host: localhost
  port: 8080

log:
  level: debug
  output: stdout`,
		}

		want := []struct {
			value     any
			retriever func(*Config) any
		}{
			{"localhost", func(c *Config) any { return c.Server.Host }},
			{8080, func(c *Config) any { return c.Server.Port }},
			{"debug", func(c *Config) any { return c.Log.Level }},
			{"stdout", func(c *Config) any { return c.Log.Output }},
		}

		cfg, err := parseConfig(mockConfig)
		if err != nil {
			t.Fatal(err)
		}

		for _, w := range want {
			if got := w.retriever(cfg); got != w.value {
				t.Errorf("Expected %v, got %v", w.value, got)
			}
		}
	})

	t.Run("testing private func loadConfig with invalid yaml", func(t *testing.T) {
		mockConfig := &mockReader{
			data: `	server:
  host: localhost
  port: 8080`,
		}

		if _, err := parseConfig(mockConfig); !strings.Contains(err.Error(), "error loading config") {
			t.Errorf("Expected error containing %q, got '%v'", "error loading config", err)
		}
	})
}

func TestFromFile(t *testing.T) {

	t.Run("testing loading config from file", func(t *testing.T) {
		fileContent := `server:
  host: localhost
  port: 8080

log:
  level: debug
  output: stdout`

		file := os.TempDir() + "/config_test.yaml"
		defer os.Remove(file)
		if err := os.WriteFile(file, []byte(fileContent), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := FromFile(file)
		if err != nil {
			t.Fatal(err)
		}

		want := []struct {
			value     any
			retriever func(*Config) any
		}{
			{"localhost", func(c *Config) any { return c.Server.Host }},
			{8080, func(c *Config) any { return c.Server.Port }},
			{"debug", func(c *Config) any { return c.Log.Level }},
			{"stdout", func(c *Config) any { return c.Log.Output }},
		}

		for _, w := range want {
			if got := w.retriever(cfg); got != w.value {
				t.Errorf("Expected %v, got %v", w.value, got)
			}
		}
	})

	t.Run("testing loading config from invalid file", func(t *testing.T) {

		file := "infalid_file"

		if _, err := FromFile(file); err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
