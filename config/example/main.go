package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/davidsugianto/go-pkgs/config"
)

// AppConfig represents the application configuration structure
type AppConfig struct {
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

type SimpleConfig struct {
	Name string `json:"name" yaml:"name"`
	Age  int    `json:"age" yaml:"age"`
}

func main() {
	fmt.Println("=== Config Loader Example ===")
	fmt.Println()

	// Example 1: Load JSON config
	fmt.Println("1. Loading JSON configuration")
	jsonConfig, err := config.LoadJSON[AppConfig]("config.json")
	if err != nil {
		log.Fatalf("Failed to load JSON config: %v", err)
	}
	printConfig("JSON Config", jsonConfig)
	fmt.Println()

	// Example 2: Load YAML config
	fmt.Println("2. Loading YAML configuration")
	yamlConfig, err := config.LoadYAML[AppConfig]("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load YAML config: %v", err)
	}
	printConfig("YAML Config", yamlConfig)
	fmt.Println()

	// Example 3: Auto-detect format (JSON)
	fmt.Println("3. Auto-detecting format (JSON)")
	autoConfigJSON, err := config.Load[AppConfig]("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	printConfig("Auto-detected JSON Config", autoConfigJSON)
	fmt.Println()

	// Example 4: Auto-detect format (YAML)
	fmt.Println("4. Auto-detecting format (YAML)")
	autoConfigYAML, err := config.Load[AppConfig]("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	printConfig("Auto-detected YAML Config", autoConfigYAML)
	fmt.Println()

	// Example 5: Load from absolute path
	fmt.Println("5. Loading from absolute path")
	absPath, _ := filepath.Abs("config.json")
	absConfig, err := config.LoadJSON[AppConfig](absPath)
	if err != nil {
		log.Fatalf("Failed to load config from absolute path: %v", err)
	}
	fmt.Printf("Loaded config from: %s\n", absPath)
	fmt.Printf("  App Name: %s\n", absConfig.AppName)
	fmt.Println()

	// Example 6: Simple config structure
	fmt.Println("6. Simple configuration structure")
	simpleConfig, err := config.LoadJSON[SimpleConfig]("simple.json")
	if err != nil {
		log.Fatalf("Failed to load simple config: %v", err)
	}
	fmt.Printf("  Name: %s\n", simpleConfig.Name)
	fmt.Printf("  Age: %d\n", simpleConfig.Age)
	fmt.Println()

	fmt.Println("=== Example Complete ===")
}

func printConfig(title string, cfg AppConfig) {
	fmt.Printf("%s:\n", title)
	fmt.Printf("  App Name: %s\n", cfg.AppName)
	fmt.Printf("  Port: %d\n", cfg.Port)
	fmt.Printf("  Debug: %v\n", cfg.Debug)
	fmt.Printf("  Database:\n")
	fmt.Printf("    Host: %s\n", cfg.Database.Host)
	fmt.Printf("    Port: %d\n", cfg.Database.Port)
	fmt.Printf("    Username: %s\n", cfg.Database.Username)
	fmt.Printf("  Endpoints: %v\n", cfg.Endpoints)
	fmt.Printf("  Metadata: %v\n", cfg.Metadata)
}
