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
    Shutdown() error
}
```

### Example Usage

```go

// Define a daemon
myService := &MyService{}

// Define another daemon
srv := NewHTTPServer(addr, handler)

if err := goblin.Run(
	goblin.WithLogFuncs(logger.Info, logger.Error),
	goblin.WithService(myService, srv),
); err != nil {
    logger.Error("goblin run", "cause", err)
}
```

Use `RunContext` to run the services with a custom `context.Context`.

```go

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

if err := goblin.RunContext(
	ctx,
	goblin.WithLogbook(logger),
	goblin.WithService(myService, srv),
); err != nil {
    logger.Error("goblin run", "cause", err)
}
```

### License

Licensed under the MIT License. See [LICENSE](./LICENSE) for more.
