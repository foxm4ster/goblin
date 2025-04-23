# Goblin: The Daemon Tamer 🧙‍♂️👹

[![GoDoc](https://godoc.org/github.com/foxm4ster/goblin?status.svg)](https://godoc.org/github.com/foxm4ster/goblin)
[![Go Report Card](https://goreportcard.com/badge/github.com/foxm4ster/goblin)](https://goreportcard.com/report/github.com/foxm4ster/goblin)
![License](https://img.shields.io/dub/l/vibe-d.svg)

![Goblin logo](https://github.com/user-attachments/assets/4cc9f068-9f31-424e-a353-2f0c645f48c8)

Welcome to **Goblin**, a deceptively simple Go library for managing daemons — small in size, mighty in magic, and always up to a bit of mischief. 🧙‍♂️✨

## Features

- 🌀 **Awaken Your Daemons**: Goblin awakens your daemons, starting them up and putting them to work.
- 🪄 **Graceful Shutdown**: Goblin ensures your daemons rest peacefully.
- 🧻 **Goblin Vibe**: Every action, failure, and success is sprinkled with a dash of goblin chaos.

## How It Works

Goblin handles multiple daemons(horde) and controls their lifecycle. Each daemon is a mischievous creature in its own right, and Goblin ensures they follow your commands (or at least tries its best!). With a playful approach to error handling and logging, Goblin never fails to entertain while it works its magic.

### Daemon Interface

In your codebase you just need to implement the `Daemon` interface to pass it into Goblin. Goblin will handle the rest.

```go
type Daemon interface {
    // The name of the daemon (used for logs and tracking)
    Name() string

    // Bring the daemon to life!
    Serve() error

    // Shutdown the daemon gracefully (hopefully without a fight)
    Shutdown() error
}
```

### Example Usage

```go

// Define a daemon
myDaemon := &MyDaemon{}

// Define another daemon
srv := NewHTTPServer(addr, handler)

if err := goblin.Awaken(
	goblin.WithLogbook(logger),
	goblin.WithDaemon(myDaemon, srv),
); err != nil {
    logger.Error("goblin couldn’t awaken", "cause", err)
}
```

Use `AwakenContext` to awaken the daemons with a custom `context.Context`. This is useful when you want to manage cancellation or timeouts more precisely.

```go

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

if err := goblin.AwakenContext(ctx,
	goblin.WithLogbook(logger),
	goblin.WithDaemon(myDaemon, srv),
); err != nil {
    logger.Error("goblin couldn’t awaken", "cause", err)
}
```

### License

Licensed under the MIT License. See [LICENSE](./LICENSE) for more.
