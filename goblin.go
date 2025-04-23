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

func Awaken(opts ...Option) error {
	return awaken(context.Background(), opts...)
}

func AwakenContext(ctx context.Context, opts ...Option) error {
	return awaken(ctx, opts...)
}

func awaken(parent context.Context, opts ...Option) error {
	manifest := &Manifest{}

	for _, opt := range opts {
		opt(manifest)
	}

	notifyCtx, cancel := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	group, ctx := errgroup.WithContext(notifyCtx)

	for _, d := range manifest.horde {
		group.Go(tinker(ctx, manifest.book, d))
	}

	if err := group.Wait(); err != nil {
		return err
	}

	if manifest.book != nil {
		manifest.book.Info("daemons asleep, goblin vanishes")
	}

	return nil
}

func tinker(ctx context.Context, book *slog.Logger, daemon Daemon) func() error {
	return func() error {
		ch := make(chan error, 1)
		defer close(ch)

		go func() {
			if book != nil {
				book.Info("goblin is tinkering with ...", "name", daemon.Name())
			}

			if err := daemon.Serve(); err != nil {
				ch <- err
			}
		}()

		select {
		case err := <-ch:
			if book != nil {
				book.Error("goblin couldn’t summon the daemon - it backfired", "name", daemon.Name(), "cause", err.Error())
			}
			return err
		case <-ctx.Done():
			if err := daemon.Shutdown(); err != nil {
				if book != nil {
					book.Error("goblin couldn’t silence the daemon - the hush spell failed", "name", daemon.Name(), "cause", err.Error())
				}
				return err
			}

			if book != nil {
				book.Info("goblin silenced the daemon, it's now resting", "name", daemon.Name())
			}

			return nil
		}
	}
}
