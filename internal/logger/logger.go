package logger

import (
	"log/slog"
	"os"
	"path/filepath"
)

// InitLogger initializes the global logger.
// It creates a log file in the user's home directory and sets it as the default logger.
func InitLogger() {
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("failed to get user home directory", "error", err)
		os.Exit(1)
	}

	logDir := filepath.Join(home, ".virak-cli", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		slog.Error("failed to create log directory", "error", err)
		os.Exit(1)
	}

	logFile, err := os.OpenFile(filepath.Join(logDir, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("failed to open log file", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(logFile, nil))
	slog.SetDefault(logger)
}
