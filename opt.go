package goblin

import (
	"log/slog"
)

type Manifest struct {
	book  *slog.Logger
	horde []Daemon
}

type Option func(*Manifest)

func WithDaemon(horde ...Daemon) Option {
	return func(m *Manifest) {
		m.horde = horde
	}
}

func WithLogbook(book *slog.Logger) Option {
	return func(m *Manifest) {
		m.book = book
	}
}
