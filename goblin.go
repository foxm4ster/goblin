package goblin

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Service interface {
	ID() string
	Serve() error
	Shutdown(context.Context) error
}

func Run(opts []Option, srvs ...Service) error {
	return run(context.Background(), opts, srvs...)
}

func RunContext(ctx context.Context, opts []Option, srvs ...Service) error {
	return run(ctx, opts, srvs...)
}

func run(parent context.Context, opts []Option, srvs ...Service) error {
	conf := Config{}

	for _, opt := range opts {
		opt(&conf)
	}

	if conf.logInfo == nil || conf.logErr == nil {
		conf.logInfo = discard
		conf.logErr = discard
	}

	notifyCtx, cancel := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	group, ctx := errgroup.WithContext(notifyCtx)

	for _, srv := range srvs {
		group.Go(handler(ctx, conf, srv))
	}

	if err := group.Wait(); err != nil {
		return err
	}

	conf.logInfo("goblin has shut down all services")

	return nil
}

func handler(ctx context.Context, conf Config, srv Service) func() error {
	return func() error {
		ch := make(chan error, 1)
		defer close(ch)

		go func() {
			conf.logInfo("goblin is starting the service", "id", srv.ID())

			if err := srv.Serve(); err != nil {
				ch <- err
			}
		}()

		select {
		case err := <-ch:
			conf.logErr("goblin couldn't start the service", "id", srv.ID(), "cause", err.Error())
			return err
		case <-ctx.Done():
			sCtx, cancel := context.WithTimeout(context.Background(), conf.shutdownTimeout)
			defer cancel()

			if conf.shutdownTimeout == 0 {
				sCtx = context.Background()
			}

			if err := srv.Shutdown(sCtx); err != nil {
				conf.logErr("goblin couldn't shut down the service", "id", srv.ID(), "cause", err.Error())
				return err
			}

			conf.logInfo("goblin successfully shut down the service", "id", srv.ID())
			return nil
		}
	}
}

func discard(msg string, args ...any) {}
