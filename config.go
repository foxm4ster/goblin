package goblin

import "time"

type LogFunc func(msg string, args ...any)

type Config struct {
	logInfo, logErr LogFunc
	shutdownTimeout time.Duration
}

func (c Config) WithLogFuncs(info, err LogFunc) Config {
	if info == nil || err == nil {
		return c
	}

	c.logInfo = info
	c.logErr = err

	return c
}

func (c Config) WithShutdownTimeout(v time.Duration) Config {
	c.shutdownTimeout = v
	return c
}
