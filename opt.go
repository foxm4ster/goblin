package goblin

type LogFunc func(msg string, args ...any)

type Config struct {
	servers []Server
	logInfo LogFunc
	logErr  LogFunc
}

type Option func(*Config)

func WithServer(s ...Server) Option {
	return func(c *Config) {
		c.servers = s
	}
}

func WithLogFuncs(info, error LogFunc) Option {
	return func(c *Config) {
		if info == nil || error == nil {
			return
		}

		c.logInfo = info
		c.logErr = error
	}
}
