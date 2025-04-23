package goblin

import (
	"log/slog"
)

type Manifest struct {
	horde []Daemon
	info func(msg string, args ...any)
	error  func(msg string, args ...any)
}

type Option func(*Manifest)

func WithDaemon(horde ...Daemon) Option {
	return func(m *Manifest) {
		m.horde = horde
	}
}

func WithLogbook(book *slog.Logger) Option {
	return func(m *Manifest) {
		if book == nil {
			return
		}

		m.info = book.Info
		m.error = book.Error
	}
}
