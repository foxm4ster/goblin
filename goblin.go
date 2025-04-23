package goblin

import (
	"context"
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

	if manifest.info == nil || manifest.error == nil {
		manifest.info = discard
		manifest.error = discard
	}

	notifyCtx, cancel := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	group, ctx := errgroup.WithContext(notifyCtx)

	for _, d := range manifest.horde {
		group.Go(tinker(ctx, d, manifest.info, manifest.error))
	}

	if err := group.Wait(); err != nil {
		return err
	}

	manifest.info("daemons asleep, goblin vanishes")

	return nil
}

func tinker(ctx context.Context, daemon Daemon, infof, errf func(msg string, args ...any)) func() error {
	return func() error {
		ch := make(chan error, 1)
		defer close(ch)

		go func() {
			infof("goblin is tinkering with ...", "name", daemon.Name())

			if err := daemon.Serve(); err != nil {
				ch <- err
			}
		}()

		select {
		case err := <-ch:
			errf("goblin couldn’t summon the daemon - it backfired", "name", daemon.Name(), "cause", err.Error())
			return err
		case <-ctx.Done():
			if err := daemon.Shutdown(); err != nil {
				errf("goblin couldn’t silence the daemon - the hush spell failed", "name", daemon.Name(), "cause", err.Error())
				return err
			}

			infof("goblin silenced the daemon, it's now resting", "name", daemon.Name())
		}

		return nil
	}
}

func discard(msg string, args ...any) {}
