package am

type Handler struct {
	Core Core
}

func NewHandler(opts ...Option) *Handler {
	core := NewCore(opts...)
	return &Handler{
		Core: core,
	}
}

func (h *Handler) Log() Logger {
	return h.Core.Log()
}

func (h *Handler) Cfg() *Config {
	return h.Core.Cfg()
}
