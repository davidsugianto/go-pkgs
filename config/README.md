# Config - Configuration Loader

A lightweight, type-safe configuration loader for Go that supports JSON and YAML formats with automatic format detection.

## Installation

```bash
go get github.com/davidsugianto/go-pkgs/config
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/davidsugianto/go-pkgs/config"
)

type AppConfig struct {
    AppName string `json:"app_name" yaml:"app_name"`
    Port    int    `json:"port" yaml:"port"`
    Debug   bool   `json:"debug" yaml:"debug"`
}

func main() {
    // Auto-detect format based on file extension
    cfg, err := config.Load[AppConfig]("config.json")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("App: %s, Port: %d\n", cfg.AppName, cfg.Port)
}
```

## Features

- ✅ **Type-Safe** - Uses Go generics for compile-time type safety
- ✅ **Multiple Formats** - Supports JSON and YAML (`.json`, `.yaml`, `.yml`)
- ✅ **Auto-Detection** - Automatically detects file format from extension
- ✅ **Simple API** - Clean, intuitive interface
- ✅ **Zero Dependencies** - Only uses standard library plus `gopkg.in/yaml.v3`

## Usage

### Define Your Config Structure

```go
type AppConfig struct {
    AppName  string            `json:"app_name" yaml:"app_name"`
    Port     int               `json:"port" yaml:"port"`
    Debug    bool              `json:"debug" yaml:"debug"`
    Database DatabaseConfig    `json:"database" yaml:"database"`
    Metadata map[string]string `json:"metadata" yaml:"metadata"`
}

type DatabaseConfig struct {
    Host     string `json:"host" yaml:"host"`
    Port     int    `json:"port" yaml:"port"`
    Username string `json:"username" yaml:"username"`
    Password string `json:"password" yaml:"password"`
}
```

**Important:** Make sure to include both `json` and `yaml` tags for maximum compatibility.

### Loading Configurations

#### Auto-Detect Format

The `Load` function automatically detects the file format based on the extension:

```go
// Detects JSON from .json extension
cfg, err := config.Load[AppConfig]("config.json")

// Detects YAML from .yaml extension
cfg, err := config.Load[AppConfig]("config.yaml")

// Detects YAML from .yml extension
cfg, err := config.Load[AppConfig]("config.yml")
```

#### Explicit Format Loading

You can also explicitly specify the format:

```go
// Load JSON explicitly
cfg, err := config.LoadJSON[AppConfig]("config.json")

// Load YAML explicitly
cfg, err := config.LoadYAML[AppConfig]("config.yaml")
```

### Example Configurations

#### JSON Example (`config.json`)

```json
{
  "app_name": "my-app",
  "port": 8080,
  "debug": true,
  "database": {
    "host": "localhost",
    "port": 5432,
    "username": "admin",
    "password": "secret123"
  },
  "metadata": {
    "version": "1.0.0",
    "environment": "development"
  }
}
```

#### YAML Example (`config.yaml`)

```yaml
app_name: my-app
port: 8080
debug: true
database:
  host: localhost
  port: 5432
  username: admin
  password: secret123
metadata:
  version: 1.0.0
  environment: development
```

### Error Handling

All functions return errors that should be checked:

```go
cfg, err := config.Load[AppConfig]("config.json")
if err != nil {
    // Handle errors:
    // - File not found
    // - Invalid JSON/YAML syntax
    // - Unsupported file format
    // - Type mismatch
    log.Fatalf("Failed to load config: %v", err)
}
```

### Supported Features

- **Nested Structures** - Full support for nested structs
- **Slices/Arrays** - Support for arrays and slices
- **Maps** - Support for map types
- **Primitives** - All Go primitive types (string, int, bool, float64, etc.)
- **Pointers** - Support for pointer types

### Complete Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/davidsugianto/go-pkgs/config"
)

type Config struct {
    AppName   string   `json:"app_name" yaml:"app_name"`
    Port      int      `json:"port" yaml:"port"`
    Endpoints []string `json:"endpoints" yaml:"endpoints"`
}

func main() {
    // Load JSON config
    cfg, err := config.Load[Config]("config.json")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("App: %s\n", cfg.AppName)
    fmt.Printf("Port: %d\n", cfg.Port)
    fmt.Printf("Endpoints: %v\n", cfg.Endpoints)
}
```

## API Reference

### Functions

#### `Load[T any](path string) (T, error)`

Automatically detects the file format based on the file extension and loads the configuration.

- `path`: Path to the configuration file (`.json`, `.yaml`, or `.yml`)
- Returns: The loaded configuration and an error

**Supported extensions:**
- `.json` - JSON format
- `.yaml` - YAML format
- `.yml` - YAML format

#### `LoadJSON[T any](path string) (T, error)`

Loads a JSON configuration file.

- `path`: Path to the JSON configuration file
- Returns: The loaded configuration and an error

#### `LoadYAML[T any](path string) (T, error)`

Loads a YAML configuration file.

- `path`: Path to the YAML configuration file (`.yaml` or `.yml`)
- Returns: The loaded configuration and an error

## Error Handling

The package returns descriptive errors for common scenarios:

- **File not found**: `failed to read config file: open <path>: no such file or directory`
- **Invalid JSON**: `failed to parse JSON config: invalid character...`
- **Invalid YAML**: `failed to parse YAML config: ...`
- **Unsupported format**: `unsupported file format: .txt (supported: .json, .yaml, .yml)`

## Best Practices

1. **Use both JSON and YAML tags**: Always include both `json` and `yaml` struct tags for maximum compatibility.

2. **Validate configurations**: After loading, consider validating your configuration values.

3. **Handle errors**: Always check and handle errors appropriately.

4. **Use meaningful defaults**: Set default values in your struct initialization or use environment-specific config files.

5. **Type safety**: Take advantage of Go generics for compile-time type safety.

## Examples

See the `example/` directory for a complete working example with sample JSON and YAML files.

```bash
cd example
go run main.go
```

## Limitations

- Does not support `.env` files (use a dedicated `.env` loader package)
- Does not support environment variable overrides (planned for future release)
- File format detection is based solely on file extension
- Does not support watching config files for changes
