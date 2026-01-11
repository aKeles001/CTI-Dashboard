package logger

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

var (
	logger  *slog.Logger
	logFile *os.File
)

// Init initializes structured logging to stdout + file
func Init(logPath string) error {
	var err error

	logDir := filepath.Join(logPath, "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logFile, err = os.OpenFile(
		filepath.Join(logDir, "scraper.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err
	}

	handler := slog.NewTextHandler(
		logFile,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	)

	logger = slog.New(handler)
	log.SetOutput(os.Stdout)

	return nil
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)

	if len(args) > 0 {
		log.Printf("[INFO] %s %v\n", msg, args)
	} else {
		log.Printf("[INFO] %s\n", msg)
	}
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
	log.Println("[ERROR]", msg)
}
