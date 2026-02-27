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

### Implementation Contract

When implementing the `Service` interface, you must honor a critical contract:

- **`Shutdown()` is a promise**: When your `Shutdown(ctx)` method returns (with or without error), it must guarantee that `Serve()` will exit shortly after. Goblin waits for the `Serve()` goroutine to complete after calling `Shutdown()`. If `Serve()` does not return, the entire shutdown process will hang indefinitely.

- **Example of correct implementation** (HTTP server):
  ```go
  func (s *MyServer) Shutdown(ctx context.Context) error {
      return s.httpServer.Shutdown(ctx)  // Guaranteed to stop Serve()
  }
  ```

- **Example of incorrect implementation** (will hang):
  ```go
  func (s *MyServer) Shutdown(ctx context.Context) error {
      s.logger.Info("shutting down")
      return nil  // Returns immediately but Serve() is still running!
  }
  ```

### Logging behavior

By default, Goblin logs simple text messages. To customize logging (e.g., disable logs or use JSON format), pass a logger instance via `goblin.WithLogger(logger)`.

### Example Usage

```go
// Simple usage with defaults
myService := &MyService{}
srv := NewHTTPServer(addr, handler)

if err := goblin.Run(myService, srv); err != nil {
    // handle error
}
```

Use `RunContext` to run the services with a custom `context.Context`.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

if err := goblin.RunContext(ctx, myService, srv); err != nil {
    // handle error
}
```

For configuration, use the `With` builder pattern:

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

if err := goblin.With(
    goblin.WithLogger(logger),
    goblin.WithShutdownTimeout(time.Second * 8),
).Run(myService, srv); err != nil {
    logger.Error("goblin run", "cause", err)
}
```

Or with custom context:

```go
if err := goblin.With(
    goblin.WithLogger(logger),
    goblin.WithShutdownTimeout(time.Second * 8),
).RunContext(ctx, myService, srv); err != nil {
    logger.Error("goblin run", "cause", err)
}
```

### License

Licensed under the MIT License. See [LICENSE](./LICENSE) for more.
