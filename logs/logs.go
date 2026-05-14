package logs

import (
	"fmt"

	"github.com/rs/zerolog"
)

// Level represents log level
type Level zerolog.Level

// Level constants
const (
	DebugLevel = Level(zerolog.DebugLevel)
	InfoLevel  = Level(zerolog.InfoLevel)
	WarnLevel  = Level(zerolog.WarnLevel)
	ErrorLevel = Level(zerolog.ErrorLevel)
	FatalLevel = Level(zerolog.FatalLevel)
)

// Debug level logs

func Debug(args ...any) {
	debugLogger.Debug().Msg(fmt.Sprint(args...))
}

func Debugf(format string, args ...any) {
	debugLogger.Debug().Msgf(format, args...)
}

func DebugWithFields(msg string, kv map[string]any) {
	evt := debugLogger.Debug()
	for k, v := range kv {
		evt = evt.Interface(k, v)
	}
	evt.Msg(msg)
}

// Info level logs

func Print(args ...any) {
	infoLogger.Info().Msg(fmt.Sprint(args...))
}

func Println(args ...any) {
	infoLogger.Info().Msg(fmt.Sprint(args...))
}

func Printf(format string, args ...any) {
	infoLogger.Info().Msgf(format, args...)
}

func Info(args ...any) {
	infoLogger.Info().Msg(fmt.Sprint(args...))
}

func Infoln(args ...any) {
	infoLogger.Info().Msg(fmt.Sprint(args...))
}

func Infof(format string, args ...any) {
	infoLogger.Info().Msgf(format, args...)
}

func InfoWithFields(msg string, kv map[string]any) {
	evt := infoLogger.Info()
	for k, v := range kv {
		evt = evt.Interface(k, v)
	}
	evt.Msg(msg)
}

// Warn level logs

func Warn(args ...any) {
	warnLogger.Warn().Msg(fmt.Sprint(args...))
}

func Warnln(args ...any) {
	warnLogger.Warn().Msg(fmt.Sprint(args...))
}

func Warnf(format string, args ...any) {
	warnLogger.Warn().Msgf(format, args...)
}

func WarnWithFields(msg string, kv map[string]any) {
	evt := warnLogger.Warn()
	for k, v := range kv {
		evt = evt.Interface(k, v)
	}
	evt.Msg(msg)
}

// Error level logs

func Error(args ...any) {
	errLogger.Error().Msg(fmt.Sprint(args...))
}

func Errorln(args ...any) {
	errLogger.Error().Msg(fmt.Sprint(args...))
}

func Errorf(format string, args ...any) {
	errLogger.Error().Msgf(format, args...)
}

func ErrorWithFields(msg string, kv map[string]any) {
	evt := errLogger.Error()
	for k, v := range kv {
		evt = evt.Interface(k, v)
	}
	evt.Msg(msg)
}

func Errors(err error) {
	if err != nil {
		errLogger.Error().Err(err).Msg("error")
	}
}

// Fatal level logs

func Fatal(args ...any) {
	fatalLogger.Fatal().Msg(fmt.Sprint(args...))
}

func Fatalln(args ...any) {
	fatalLogger.Fatal().Msg(fmt.Sprint(args...))
}

func Fatalf(format string, args ...any) {
	fatalLogger.Fatal().Msgf(format, args...)
}

func FatalWithFields(msg string, kv map[string]any) {
	evt := fatalLogger.Fatal()
	for k, v := range kv {
		evt = evt.Interface(k, v)
	}
	evt.Msg(msg)
}