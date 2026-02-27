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

func Run(srvs ...Service) error {
	return run(context.Background(), srvs)
}

func RunContext(ctx context.Context, srvs ...Service) error {
	return run(ctx, srvs)
}

func run(parent context.Context, srvs []Service, opts ...Option) error {
	conf := newDefaultConfig()

	for _, opt := range opts {
		opt(&conf)
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

	conf.logger.Info("all services shut down")

	return nil
}

func handler(ctx context.Context, conf Config, srv Service) func() error {
	return func() error {
		ch := make(chan error, 1)
		done := make(chan struct{})

		go func() {
			conf.logger.Info("starting service", "id", srv.ID())

			ch <- srv.Serve()
			close(done)
			close(ch)
		}()

		select {
		case err := <-ch:
			if err != nil {
				conf.logger.Error("couldn't start service",
					"id", srv.ID(),
					"cause", err.Error())
				return err
			}

			conf.logger.Info("service terminated naturally without signal", "id", srv.ID())
			return nil
		case <-ctx.Done():
			sdCtx, cancel := context.WithTimeout(context.Background(), conf.shutdownTimeout)
			defer cancel()

			err := srv.Shutdown(sdCtx)

			<-done

			if err != nil {
				conf.logger.Error("couldn't shut down service",
					"id", srv.ID(),
					"cause", err.Error())
				return err
			}

			conf.logger.Info("successfully shut down service", "id", srv.ID())
			return nil
		}
	}
}
