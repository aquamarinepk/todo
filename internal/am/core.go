package am

type Core interface {
	Log() Logger
}

type BaseCore struct {
	log Logger
}

func NewCore(log Logger) *BaseCore {
	return &BaseCore{
		log: log,
	}
}

func (c *BaseCore) Log() Logger {
	return c.log
}
