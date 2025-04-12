package goblin

import (
	"log/slog"
)

type Config struct {
	book  *slog.Logger
	horde []Daemon
}

type Option func(*Config)

func WithDaemon(horde ...Daemon) func(*Config) {
	return func(c *Config) {
		c.horde = horde
	}
}

func WithLogbook(book *slog.Logger) func(*Config) {
	return func(c *Config) {
		if book != nil && c.book == nil {
			c.book = book
		}
	}
}
