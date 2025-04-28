package goblin

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server interface {
	ID() string
	Serve() error
	Shutdown() error
}

func Run(opts ...Option) error {
	return run(context.Background(), opts...)
}

func RunContext(ctx context.Context, opts ...Option) error {
	return run(ctx, opts...)
}

func run(parent context.Context, opts ...Option) error {
	conf := &Config{}

	for _, opt := range opts {
		opt(conf)
	}

	if conf.logInfo == nil || conf.logErr == nil {
		conf.logInfo = discard
		conf.logErr = discard
	}

	notifyCtx, cancel := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	group, ctx := errgroup.WithContext(notifyCtx)

	for _, srv := range conf.servers {
		group.Go(handler(ctx, srv, conf.logInfo, conf.logErr))
	}

	if err := group.Wait(); err != nil {
		return err
	}

	conf.logInfo("goblin has shut down all servers")

	return nil
}

func handler(ctx context.Context, srv Server, logInfo, logErr LogFunc) func() error {
	return func() error {
		ch := make(chan error, 1)
		defer close(ch)

		go func() {
			logInfo("goblin is starting the server", "id", srv.ID())

			if err := srv.Serve(); err != nil {
				ch <- err
			}
		}()

		select {
		case err := <-ch:
			logErr("goblin could't start the server", "id", srv.ID(), "cause", err.Error())
			return err
		case <-ctx.Done():
			if err := srv.Shutdown(); err != nil {
				logErr("goblin couldn't shut down the server", "id", srv.ID(), "cause", err.Error())
				return err
			}

			logInfo("goblin successfully shut down the server", "id", srv.ID())
		}

		return nil
	}
}

func discard(msg string, args ...any) {}
