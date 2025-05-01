package goblin

type LogFunc func(msg string, args ...any)

type Config struct {
	services        []Service
	logInfo, logErr LogFunc
}

type Option func(*Config)

func WithService(s ...Service) Option {
	return func(c *Config) {
		c.services = s
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
