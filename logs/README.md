# logs

Unified logging package built on top of [logger](../logger) with a simple, package-level API for structured logging.

## Quick Start

```go
import log "github.com/davidsugianto/go-pkgs/logs"

func main() {
    log.Info("application started")
    log.Debugf("processing item: %d", 42)
    log.WarnWithFields("rate limit approaching", map[string]any{
        "remaining": 5,
        "limit":     100,
    })
}
```

## Log Levels

| Function          | Level |
|-------------------|-------|
| `Debug`, `Debugf`, `DebugWithFields` | Debug |
| `Info`, `Infof`, `InfoWithFields`    | Info  |
| `Print`, `Printf`, `Println`         | Info  |
| `Warn`, `Warnf`, `WarnWithFields`    | Warn  |
| `Error`, `Errorf`, `ErrorWithFields` | Error |
| `Errors(err)`                         | Error |
| `Fatal`, `Fatalf`, `FatalWithFields` | Fatal |

## Configuration

### SetConfig

Replace all loggers with custom configuration at once.

```go
err := log.SetConfig(&log.Config{
    Level:   log.InfoLevel,
    AppName: "my-service",
    UseJSON: true,
})
```

### SetLevel

Adjust the level threshold for all loggers at runtime.

```go
log.SetLevel(log.WarnLevel)  // only Warn and above
```

### SetLogger

Replace the logger for a specific level — useful for routing errors to a separate file.

```go
errorLogger, _ := log.NewLogger(&log.Config{
    Level:   log.ErrorLevel,
    LogFile: "/var/log/app-error.log",
    UseJSON: true,
})
log.SetLogger(log.ErrorLevel, errorLogger)
```

## Config Fields

| Field      | Type    | Description                         |
|------------|---------|-------------------------------------|
| `Level`    | `Level` | Log level threshold                 |
| `AppName`  | `string`| Application name added to each log  |
| `Caller`   | `bool`  | Include caller line numbers         |
| `LogFile`  | `string`| Output file for standard logs       |
| `DebugFile`| `string`| Separate output file for debug logs |
| `UseColor` | `bool`  | Colorized console output            |
| `UseJSON`  | `bool`  | JSON format output                  |