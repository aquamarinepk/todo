package am

import "context"

type Handler struct {
	Core Core
}

func NewHandler(name string, opts ...Option) *Handler {
	core := NewCore(name, opts...)
	return &Handler{
		Core: core,
	}
}

func (h *Handler) Name() string {
	return h.Core.Name()
}

func (h *Handler) SetName(name string) {
	h.Core.SetName(name)
}

func (h *Handler) Log() Logger {
	return h.Core.Log()
}

func (h *Handler) SetLog(log Logger) {
	h.Core.SetLog(log)
}

func (h *Handler) Cfg() *Config {
	return h.Core.Cfg()
}

func (h *Handler) SetCfg(cfg *Config) {
	h.Core.SetCfg(cfg)
}

func (h *Handler) Setup(ctx context.Context) error {
	return h.Core.Setup(ctx)
}

func (h *Handler) Start(ctx context.Context) error {
	return h.Core.Start(ctx)
}

func (h *Handler) Stop(ctx context.Context) error {
	return h.Core.Stop(ctx)
}
