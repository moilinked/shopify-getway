package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

var globalFileWriter *WeeklyRotateWriter

const defaultLogDir = "./logs"

func Init(logLevel string) error {
	level, err := parseLevel(logLevel)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339

	if err := os.MkdirAll(defaultLogDir, 0755); err != nil {
		return fmt.Errorf("create log dir %s: %w", defaultLogDir, err)
	}

	globalFileWriter = &WeeklyRotateWriter{Dir: defaultLogDir, Prefix: "app"}
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}

	Log = zerolog.New(io.MultiWriter(consoleWriter, globalFileWriter)).With().Timestamp().Caller().Logger()
	return nil
}

func Close() {
	if globalFileWriter != nil {
		_ = globalFileWriter.Close()
	}
}

func parseLevel(s string) (zerolog.Level, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "trace":
		return zerolog.TraceLevel, nil
	case "debug":
		return zerolog.DebugLevel, nil
	case "info", "":
		return zerolog.InfoLevel, nil
	case "warn", "warning":
		return zerolog.WarnLevel, nil
	case "error":
		return zerolog.ErrorLevel, nil
	case "fatal":
		return zerolog.FatalLevel, nil
	default:
		return zerolog.InfoLevel, fmt.Errorf("unknown log level: %s", s)
	}
}
