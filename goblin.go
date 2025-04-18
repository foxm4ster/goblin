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

type Goblin struct {
	horde  []Daemon
	scrawl Scrawler
}

func New(opts ...Option) Goblin {
	man := &Manifest{}

	for _, opt := range opts {
		opt(man)
	}

	return Goblin{
		scrawl: withLogbook(man.book),
		horde:  man.horde,
	}
}

func (g Goblin) Awaken() error {
	return g.awaken(context.Background())
}

func (g Goblin) AwakenContext(ctx context.Context) error {
	return g.awaken(ctx)
}

func (g Goblin) awaken(parent context.Context) error {
	notifyCtx, cancel := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	group, ctx := errgroup.WithContext(notifyCtx)

	for _, d := range g.horde {
		group.Go(tinker(ctx, g.scrawl, d))
	}

	if err := group.Wait(); err != nil {
		return err
	}

	g.scrawl(slog.LevelInfo, "daemons asleep, goblin vanishes")

	return nil
}

func tinker(ctx context.Context, scrawl Scrawler, daemon Daemon) func() error {
	return func() error {
		ch := make(chan error, 1)
		defer close(ch)

		go func() {
			scrawl(slog.LevelInfo, "goblin is tinkering with ...", slog.String("name", daemon.Name()))

			if err := daemon.Serve(); err != nil {
				ch <- err
			}
		}()

		select {
		case err := <-ch:
			scrawl(
				slog.LevelError,
				"goblin couldn’t summon the daemon",
				slog.String("name", daemon.Name()),
				slog.String("cause", err.Error()),
			)
			return err
		case <-ctx.Done():
			if err := daemon.Shutdown(); err != nil {
				scrawl(
					slog.LevelError,
					"goblin couldn’t silence the daemon",
					slog.String("name", daemon.Name()),
					slog.String("cause", err.Error()),
				)
				return err
			}

			scrawl(slog.LevelInfo, "goblin silenced the daemon", slog.String("name", daemon.Name()))
			return nil
		}
	}
}
