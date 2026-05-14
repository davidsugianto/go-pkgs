package logs

import (
	"errors"
	"fmt"
	"os"

	"github.com/davidsugianto/go-pkgs/logger"
	"github.com/rs/zerolog"
)

// Config of log
type Config struct {
	// Level log level, default: DebugLevel
	Level Level

	// AppName is the application name this log belong to
	AppName string

	// Caller, option to print caller line numbers.
	Caller bool

	// LogFile is output file for log other than debug log
	LogFile string

	// DebugFile is output file for debug log
	DebugFile string

	// UseColor, option to colorize log in console.
	UseColor bool

	// UseJSON, option to print in json format.
	UseJSON bool
}

var (
	infoLogger  = newDefault(zerolog.InfoLevel)
	debugLogger = newDefault(zerolog.DebugLevel)
	warnLogger  = infoLogger
	errLogger   = infoLogger
	fatalLogger = infoLogger
)

// newDefault creates a default logger for the given level
func newDefault(level zerolog.Level) *logger.Logger {
	return logger.NewWithConfig(logger.Config{
		Level:  level,
		Format: logger.FormatConsole,
	})
}

// NewLogger creates a new logger with the given configuration
func NewLogger(config *Config) (*logger.Logger, error) {
	if config == nil {
		config = &Config{Level: InfoLevel}
	}

	cfg := logger.Config{
		Level:       zerolog.Level(config.Level),
		ServiceName: config.AppName,
		Caller:      config.Caller,
	}

	// Format
	if config.UseJSON {
		cfg.Format = logger.FormatJSON
	} else if config.UseColor {
		cfg.Format = logger.FormatPretty
	} else {
		cfg.Format = logger.FormatConsole
	}

	// Output to file if specified
	if config.LogFile != "" {
		f, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file %s: %w", config.LogFile, err)
		}
		cfg.Output = f
	}

	return logger.NewWithConfig(cfg), nil
}

// SetLogger sets the logger for a specific level
func SetLogger(level Level, l *logger.Logger) error {
	if level < DebugLevel || level > FatalLevel {
		return errors.New("invalid level")
	}
	if l == nil {
		return errors.New("invalid logger")
	}
	*loggers[level] = l
	return nil
}

// loggers index maps levels to their logger pointers
var loggers = [5]**logger.Logger{
	&debugLogger,
	&infoLogger,
	&warnLogger,
	&errLogger,
	&fatalLogger,
}

// SetLevel adjusts log level threshold for all loggers
func SetLevel(level Level) {
	if level < DebugLevel {
		level = InfoLevel
	}

	debugLogger.SetLevel(zerolog.Level(level))
	infoLogger.SetLevel(zerolog.Level(level))
	warnLogger.SetLevel(zerolog.Level(level))
	errLogger.SetLevel(zerolog.Level(level))
	fatalLogger.SetLevel(zerolog.Level(level))
}

// SetConfig creates new default (info & debug) loggers based on given config
func SetConfig(config *Config) error {
	var debugCfg, infoCfg Config
	if config != nil {
		infoCfg = *config
		infoCfg.Level = InfoLevel

		debugCfg = *config
		debugCfg.Level = DebugLevel
		if debugCfg.DebugFile != "" {
			debugCfg.LogFile = debugCfg.DebugFile
		}
	}

	newLogger, err := NewLogger(&infoCfg)
	if err != nil {
		return err
	}
	infoLogger = newLogger
	warnLogger = newLogger
	errLogger = newLogger
	fatalLogger = newLogger

	newDebugLogger, err := NewLogger(&debugCfg)
	if err != nil {
		return err
	}
	debugLogger = newDebugLogger

	return nil
}