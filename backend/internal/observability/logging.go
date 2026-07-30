package observability

import (
	"fmt"
	"io"
	"log/slog"
	"qvarkk/kv2/internal/buildinfo"
)

type LogConfig struct {
	Format    string // "json", "text"
	Level     string // "debug", "info", "warn", "error"
	AddSource bool
}

type Logging struct {
	Logger *slog.Logger
	Level  *slog.LevelVar
}

func NewLogger(
	output io.Writer,
	cfg LogConfig,
	buildInfo buildinfo.Info,
) (*Logging, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	levelVar := new(slog.LevelVar)
	levelVar.Set(level)

	options := &slog.HandlerOptions{
		Level:     levelVar,
		AddSource: cfg.AddSource,
	}

	var handler slog.Handler

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(output, options)
	case "text":
		handler = slog.NewTextHandler(output, options)
	default:
		return nil, fmt.Errorf("unsupported log format %q", cfg.Format)
	}

	logger := slog.New(handler).With(
		slog.String("version", buildInfo.Version),
		slog.String("go_version", buildInfo.GoVersion),
	)

	return &Logging{
		Logger: logger,
		Level:  levelVar,
	}, nil
}

func parseLevel(value string) (slog.Level, error) {
	switch value {
	case "debug":
		return slog.LevelDebug, nil
	case "", "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("unsupported log level %q", value)
	}
}
