package goblin

import (
	"log/slog"
)

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
