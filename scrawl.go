package goblin

import (
	"context"
	"log/slog"
)

type Scrawler func(level slog.Level, msg string, attr ...slog.Attr)

func withLogbook(book *slog.Logger) Scrawler {
	return func(level slog.Level, msg string, attr ...slog.Attr) {
		if book == nil {
			return
		}

		book.LogAttrs(context.Background(), level, msg, attr...)
	}
}
