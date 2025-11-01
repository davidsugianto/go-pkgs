package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func Load[T any](path string) (T, error) {
	var config T

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return LoadJSON[T](path)
	case ".yaml", ".yml":
		return LoadYAML[T](path)
	default:
		return config, fmt.Errorf("unsupported file format: %s (supported: .json, .yaml, .yml)", ext)
	}
}

func LoadJSON[T any](path string) (T, error) {
	var config T

	data, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse JSON config: %w", err)
	}

	return config, nil
}

func LoadYAML[T any](path string) (T, error) {
	var config T

	data, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	return config, nil
}
