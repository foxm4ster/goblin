package goblin

import (
	"io"
	"log/slog"
	"time"
)

type Option func(*Config)

func WithLogger(logger *slog.Logger) Option {
	return func(c *Config) {
		if logger == nil {
			return
		}
		c.logger = logger
	}
}

func WithNopLogger() Option {
	return func(c *Config) {
		c.logger = slog.New(
			slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}),
		)
	}
}

func WithShutdownTimeout(t time.Duration) Option {
	return func(c *Config) {
		if t > 0 {
			c.shutdownTimeout = t
		}
	}
}
