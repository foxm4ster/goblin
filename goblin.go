package goblin

import (
	"context"
	"io"
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
	horde []Daemon
	book  *slog.Logger
}

func New(opts ...Option) Goblin {
	man := &Manifest{
		book: slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
	}

	for _, opt := range opts {
		opt(man)
	}

	return Goblin{
		book:  man.book,
		horde: man.horde,
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
		group.Go(tinker(ctx, g.book, d))
	}

	if err := group.Wait(); err != nil {
		return err
	}

	g.book.Info("daemons asleep, goblin vanishes")

	return nil
}

func tinker(ctx context.Context, book *slog.Logger, daemon Daemon) func() error {
	return func() error {
		ch := make(chan error, 1)
		defer close(ch)

		go func() {
			book.Info("goblin is tinkering with ...", "name", daemon.Name())

			if err := daemon.Serve(); err != nil {
				ch <- err
			}
		}()

		select {
		case err := <-ch:
			book.Error("goblin couldn’t summon the daemon", "name", daemon.Name(), "cause", err.Error())
			return err
		case <-ctx.Done():
			if err := daemon.Shutdown(); err != nil {
				book.Error("goblin couldn’t silence the daemon", "name", daemon.Name(), "cause", err.Error())
				return err
			}

			book.Info("goblin silenced the daemon", "name", daemon.Name())
			return nil
		}
	}
}
