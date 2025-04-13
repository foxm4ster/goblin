# Goblin: The Daemon Tamer 🧙‍♂️👹

Welcome to **Goblin**, a mischievous Go library for managing daemons like a true goblin would — with magic, chaos, and just a touch of mayhem. 🧙‍♂️✨

## Features

- 🌀 **Tame Your Daemons**: Goblin can awaken, silence, and tame your daemons.
- 🪄 **Graceful Shutdown**: Goblin ensures your daemons rest peacefully.
- 🧻 **Goblin Vibe**: Every action, failure, and success is sprinkled with a dash of goblin chaos.

## How It Works

Goblin creates a **daemon manager** that can handle multiple daemons(horde) and control their lifecycle. Each daemon is a mischievous creature in its own right, and Goblin ensures they follow your commands (or at least tries its best!). With a playful approach to error handling and logging, Goblin never fails to entertain while it works its magic.

### Daemon Interface

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

gob := goblin.New(
    goblin.WithLogbook(logger),
    goblin.WithDaemon(myDaemon),
)

if err := gob.Awaken(); err != nil {
    logger.Error("goblin couldn’t tame the daemons", "cause", err)
}
```


### License

Licensed under the MIT License. See LICENSE for more.
