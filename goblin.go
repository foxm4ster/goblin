package goblin

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Daemon interface {
	Name() string
	Serve() error
	Shutdown() error
}

type Config struct {
	book  *slog.Logger
	horde []Daemon
}

type Option func(*Config)

type Goblin struct {
	horde    []Daemon
	scrawler Scrawler
}

func New(opts ...Option) Goblin {
	conf := &Config{}

	for _, opt := range opts {
		opt(conf)
	}

	return Goblin{
		scrawler: withLogbook(conf.book),
		horde:    conf.horde,
	}
}

func (g Goblin) Awaken() error {
	parent, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	group, ctx := errgroup.WithContext(parent)

	for _, d := range g.horde {
		group.Go(tinker(ctx, g.scrawler, d))
	}

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

func tinker(ctx context.Context, scrawl Scrawler, daemon Daemon) func() error {
	return func() error {
		ch := make(chan error, 1)
		defer close(ch)

		go func() {
			scrawl(slog.LevelInfo, "running service", slog.String("name", daemon.Name()))

			if err := daemon.Serve(); err != nil {
				ch <- err
			}
		}()

		select {
		case err := <-ch:
			scrawl(
				slog.LevelError,
				"failed to run service",
				slog.String("name", daemon.Name()),
				slog.String("cause", err.Error()),
			)
			return err
		case <-ctx.Done():
			if err := daemon.Shutdown(); err != nil {
				scrawl(
					slog.LevelError,
					"failed to stop service",
					slog.String("name", daemon.Name()),
					slog.String("cause", err.Error()),
				)
				return err
			}

			scrawl(slog.LevelInfo, "service is stopping", slog.String("name", daemon.Name()))
			return nil
		}
	}
}
