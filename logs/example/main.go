package main

import (
	"os"

	log "github.com/davidsugianto/go-pkgs/logs"
)

func main() {
	// Example 1: Quick usage with defaults (console output, info level)
	log.Info("application started")
	log.Debugf("debug message with value: %d", 42)
	log.WarnWithFields("rate limit warning", map[string]any{
		"current_rate": 95,
		"limit":        100,
	})

	// Example 2: SetLevel to control what gets printed
	log.SetLevel(log.WarnLevel)
	log.Info("this won't be printed")   // filtered out
	log.Debug("this won't be printed")  // filtered out
	log.Warn("this will be printed")    // printed

	// Reset level
	log.SetLevel(log.DebugLevel)

	// Example 3: Configure with SetConfig
	err := log.SetConfig(&log.Config{
		Level:   log.InfoLevel,
		AppName: "my-service",
		UseJSON: true,
	})
	if err != nil {
		log.Fatalf("failed to configure logger: %v", err)
	}
	log.InfoWithFields("service initialized", map[string]any{
		"version": "1.0.0",
		"port":    8080,
	})

	// Example 4: Per-level logger with file output
	debugLogger, err := log.NewLogger(&log.Config{
		Level:   log.DebugLevel,
		AppName: "my-service",
		LogFile: "/tmp/app.log",
		UseJSON: true,
	})
	if err != nil {
		log.Fatalln(err)
	}
	_ = debugLogger

	// Example 5: SetLogger for specific level
	errorFileLogger, err := log.NewLogger(&log.Config{
		Level:   log.ErrorLevel,
		AppName: "my-service",
		LogFile: "/tmp/error.log",
		UseJSON: true,
	})
	if err != nil {
		log.Fatalln(err)
	}

	if err := log.SetLogger(log.ErrorLevel, errorFileLogger); err != nil {
		log.Fatalln(err)
	}

	// Example 6: Error handling
	log.Errorf("failed to process request: %v", os.ErrNotExist)
	log.ErrorWithFields("database error", map[string]any{
		"operation": "INSERT",
		"table":     "users",
		"error":     "connection timeout",
	})
	log.Errors(os.ErrPermission)

	// Example 7: Print-style compatibility
	log.Print("standard print style")
	log.Printf("formatted print: %s", "hello")
	log.Println("println style")

	// Example 8: Fatal (will exit the program - commented out)
	// log.Fatal("critical failure")
	// log.Fatalf("cannot start on port %d", 8080)
}