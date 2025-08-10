# goblin - service manager

[![GoDoc](https://godoc.org/github.com/foxm4ster/goblin?status.svg)](https://godoc.org/github.com/foxm4ster/goblin)
[![Go Report Card](https://goreportcard.com/badge/github.com/foxm4ster/goblin)](https://goreportcard.com/report/github.com/foxm4ster/goblin)
![License](https://img.shields.io/dub/l/vibe-d.svg)

`goblin` is a lightweight, powerful service management tool built for simplicity and efficiency.

### Service Interface

In your codebase you just need to implement the `Service` interface to pass it into Goblin. Goblin will handle the rest.

```go
type Service interface {
    ID() string
    Serve() error
    Shutdown(ctx context.Context) error
}
```

### Example Usage

```go

// Define a service
myService := &MyService{}

// Define another service
srv := NewHTTPServer(addr, handler)

opts := []goblin.Option{
    goblin.WithLogFuncs(logger.Info, logger.Error),
    goblin.WithShutdownTimeout(time.Second * 8),
}

if err := goblin.Run(opts, myService, srv); err != nil {
    logger.Error("goblin run", "cause", err)
}
```

Use `RunContext` to run the services with a custom `context.Context`.

```go

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

opts := []goblin.Option{
    goblin.WithLogFuncs(logger.Info, logger.Error),
    goblin.WithShutdownTimeout(time.Second * 8),
}

if err := goblin.RunContext(ctx, opts, myService, srv); err != nil {
    logger.Error("goblin run", "cause", err)
}
```

If you don't need logging or any configuration, you can pass nil.

```go
if err := goblin.Run(nil, myService, srv); err != nil {
    logger.Error("goblin run", "cause", err)
}
```

---

If you prefer a builder style, you can do it using `With`.

```go

if err := goblin.With(
    goblin.WithLogFuncs(logger.Info, logger.Error),
    goblin.WithShutdownTimeout(time.Second * 8),
).RunContext(ctx, myService, srv); err != nil {
    logger.Error("goblin run", "cause", err)
}

```

### License

Licensed under the MIT License. See [LICENSE](./LICENSE) for more.
