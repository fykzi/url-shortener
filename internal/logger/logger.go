package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func SetupLogger(logLevel string) (*slog.Logger, error) {
	switch logLevel {
	case "INFO":
		return slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		), nil
	case "WARN":
		return slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
		), nil
	case "ERROR":
		return slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}),
		), nil
	case "DEBUG":
		return slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		), nil
	default:
		return nil, fmt.Errorf("invalid logger level")
	}
}
