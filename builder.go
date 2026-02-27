package goblin

import "context"

type Builder struct {
	opts []Option
}

func With(opts ...Option) Builder {
	return Builder{opts: opts}
}

func (b Builder) Run(srvs ...Service) error {
	return run(context.Background(), srvs, b.opts...)
}

func (b Builder) RunContext(ctx context.Context, srvs ...Service) error {
	return run(ctx, srvs, b.opts...)
}
