package config

import (
	"os"
	"path/filepath"
	"testing"
)

type TestConfig struct {
	AppName   string            `json:"app_name" yaml:"app_name"`
	Port      int               `json:"port" yaml:"port"`
	Debug     bool              `json:"debug" yaml:"debug"`
	Database  DatabaseConfig    `json:"database" yaml:"database"`
	Endpoints []string          `json:"endpoints" yaml:"endpoints"`
	Metadata  map[string]string `json:"metadata" yaml:"metadata"`
}

type DatabaseConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func TestLoadJSON(t *testing.T) {
	// Create temporary JSON file
	jsonContent := `{
		"app_name": "test-app",
		"port": 8080,
		"debug": true,
		"database": {
			"host": "localhost",
			"port": 5432,
			"username": "admin",
			"password": "secret"
		},
		"endpoints": ["/api/v1", "/api/v2"],
		"metadata": {
			"version": "1.0.0",
			"environment": "test"
		}
	}`

	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(jsonFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to create test JSON file: %v", err)
	}

	config, err := LoadJSON[TestConfig](jsonFile)
	if err != nil {
		t.Fatalf("LoadJSON failed: %v", err)
	}

	if config.AppName != "test-app" {
		t.Errorf("Expected AppName 'test-app', got '%s'", config.AppName)
	}
	if config.Port != 8080 {
		t.Errorf("Expected Port 8080, got %d", config.Port)
	}
	if !config.Debug {
		t.Error("Expected Debug to be true")
	}
	if config.Database.Host != "localhost" {
		t.Errorf("Expected Database.Host 'localhost', got '%s'", config.Database.Host)
	}
	if config.Database.Port != 5432 {
		t.Errorf("Expected Database.Port 5432, got %d", config.Database.Port)
	}
	if len(config.Endpoints) != 2 {
		t.Errorf("Expected 2 endpoints, got %d", len(config.Endpoints))
	}
	if config.Metadata["version"] != "1.0.0" {
		t.Errorf("Expected metadata version '1.0.0', got '%s'", config.Metadata["version"])
	}
}

func TestLoadYAML(t *testing.T) {
	// Create temporary YAML file
	yamlContent := `app_name: test-app
port: 8080
debug: true
database:
  host: localhost
  port: 5432
  username: admin
  password: secret
endpoints:
  - /api/v1
  - /api/v2
metadata:
  version: 1.0.0
  environment: test
`

	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	config, err := LoadYAML[TestConfig](yamlFile)
	if err != nil {
		t.Fatalf("LoadYAML failed: %v", err)
	}

	if config.AppName != "test-app" {
		t.Errorf("Expected AppName 'test-app', got '%s'", config.AppName)
	}
	if config.Port != 8080 {
		t.Errorf("Expected Port 8080, got %d", config.Port)
	}
	if !config.Debug {
		t.Error("Expected Debug to be true")
	}
	if config.Database.Host != "localhost" {
		t.Errorf("Expected Database.Host 'localhost', got '%s'", config.Database.Host)
	}
	if len(config.Endpoints) != 2 {
		t.Errorf("Expected 2 endpoints, got %d", len(config.Endpoints))
	}
}

func TestLoadYML(t *testing.T) {
	// Test .yml extension
	yamlContent := `app_name: test-app
port: 8080
`

	tmpDir := t.TempDir()
	ymlFile := filepath.Join(tmpDir, "config.yml")
	if err := os.WriteFile(ymlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	config, err := LoadYAML[TestConfig](ymlFile)
	if err != nil {
		t.Fatalf("LoadYAML failed: %v", err)
	}

	if config.AppName != "test-app" {
		t.Errorf("Expected AppName 'test-app', got '%s'", config.AppName)
	}
}

func TestLoadAutoDetectJSON(t *testing.T) {
	jsonContent := `{
		"app_name": "test-app",
		"port": 8080
	}`

	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(jsonFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to create test JSON file: %v", err)
	}

	config, err := Load[TestConfig](jsonFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if config.AppName != "test-app" {
		t.Errorf("Expected AppName 'test-app', got '%s'", config.AppName)
	}
}

func TestLoadAutoDetectYAML(t *testing.T) {
	yamlContent := `app_name: test-app
port: 8080
`

	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	config, err := Load[TestConfig](yamlFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if config.AppName != "test-app" {
		t.Errorf("Expected AppName 'test-app', got '%s'", config.AppName)
	}
}

func TestLoadAutoDetectYML(t *testing.T) {
	yamlContent := `app_name: test-app
port: 8080
`

	tmpDir := t.TempDir()
	ymlFile := filepath.Join(tmpDir, "config.yml")
	if err := os.WriteFile(ymlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	config, err := Load[TestConfig](ymlFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if config.AppName != "test-app" {
		t.Errorf("Expected AppName 'test-app', got '%s'", config.AppName)
	}
}

func TestLoadUnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	txtFile := filepath.Join(tmpDir, "config.txt")
	if err := os.WriteFile(txtFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := Load[TestConfig](txtFile)
	if err == nil {
		t.Fatal("Expected error for unsupported format, got nil")
	}

	if err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}

func TestLoadJSONFileNotFound(t *testing.T) {
	_, err := LoadJSON[TestConfig]("nonexistent.json")
	if err == nil {
		t.Fatal("Expected error for non-existent file, got nil")
	}
}

func TestLoadYAMLFileNotFound(t *testing.T) {
	_, err := LoadYAML[TestConfig]("nonexistent.yaml")
	if err == nil {
		t.Fatal("Expected error for non-existent file, got nil")
	}
}

func TestLoadJSONInvalidJSON(t *testing.T) {
	invalidJSON := `{
		"app_name": "test-app"
		"port": 8080  // missing comma
	}`

	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(jsonFile, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("Failed to create test JSON file: %v", err)
	}

	_, err := LoadJSON[TestConfig](jsonFile)
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

func TestLoadYAMLInvalidYAML(t *testing.T) {
	invalidYAML := `app_name: test-app
port: 8080
  invalid_indentation: wrong
`

	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(yamlFile, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	_, err := LoadYAML[TestConfig](yamlFile)
	if err == nil {
		t.Fatal("Expected error for invalid YAML, got nil")
	}
}

func TestLoadWithEmptyStruct(t *testing.T) {
	type EmptyConfig struct{}

	jsonContent := `{}`

	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(jsonFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to create test JSON file: %v", err)
	}

	config, err := LoadJSON[EmptyConfig](jsonFile)
	if err != nil {
		t.Fatalf("LoadJSON failed: %v", err)
	}

	// Verify config was loaded (EmptyConfig is a struct, so it's always non-nil)
	_ = config
}
