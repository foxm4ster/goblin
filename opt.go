package goblin

import "time"

type LogFunc func(msg string, args ...any)

type Config struct {
	logInfo, logErr LogFunc
	shutdownTimeout time.Duration
}

type Option func(*Config)

func WithLogFuncs(info, err LogFunc) Option {
	return func(c *Config) {
		if info == nil || err == nil {
			return
		}

		c.logInfo = info
		c.logErr = err
	}
}

func WithShutdownTimeout(v time.Duration) Option {
	return func(c *Config) {
		c.shutdownTimeout = v
	}
}
