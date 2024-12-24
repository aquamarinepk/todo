package am

type Handler struct {
	Core Core
}

func NewHandler(log Logger) *Handler {
	core := NewCore(log)
	return &Handler{
		Core: core,
	}
}

func (h *Handler) Log() Logger {
	return h.Core.Log()
}
