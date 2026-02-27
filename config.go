package goblin

import (
	"log/slog"
	"os"
	"time"
)

type Config struct {
	logger          *slog.Logger
	shutdownTimeout time.Duration
}

func newDefaultConfig() Config {
	return Config{
		logger: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}),
		),
		shutdownTimeout: time.Second * 30,
	}
}
