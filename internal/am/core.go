package am

type Core interface {
	Log() Logger
	Cfg() *Config
}

type BaseCore struct {
	log Logger
	cfg *Config
}

// NewCore creates a new BaseCore instance with the provided opts.
func NewCore(opts ...Option) *BaseCore {
	core := &BaseCore{}
	for _, opt := range opts {
		opt(core)
	}
	return core
}

func (c *BaseCore) Log() Logger {
	return c.log
}

func (c *BaseCore) Cfg() *Config {
	return c.cfg
}

// Option defines a type for setting optional parameters in BaseCore.
type Option func(*BaseCore)

// WithLog sets the Logger in BaseCore.
func WithLog(log Logger) Option {
	return func(c *BaseCore) {
		c.log = log
	}
}

// WithCfg sets the Config in BaseCore.
func WithCfg(cfg *Config) Option {
	return func(c *BaseCore) {
		c.cfg = cfg
	}
}
