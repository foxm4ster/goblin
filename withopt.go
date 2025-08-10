package goblin

import "context"

type WithOpt struct {
	opts []Option
}

func With(opts ...Option) WithOpt {
	return WithOpt{opts: opts}
}

func (w WithOpt) Run(srvs ...Service) error {
	return Run(w.opts, srvs...)
}

func (w WithOpt) RunContext(ctx context.Context, srvs ...Service) error {
	return RunContext(ctx, w.opts, srvs...)
}
